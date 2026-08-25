import json
import unittest
import xml.etree.ElementTree as element_tree
from pathlib import Path
from typing import Final


ROOT: Final = Path(__file__).resolve().parents[1]
SKILL_PATH: Final = ROOT / "skills" / "next" / "SKILL.md"
OPENAI_PATH: Final = ROOT / "skills" / "next" / "agents" / "openai.yaml"
MANIFEST_PATH: Final = ROOT / ".codex-plugin" / "plugin.json"


def parse_frontmatter(document: str) -> tuple[dict[str, str], str]:
    lines = document.splitlines()
    if not lines or lines[0] != "---":
        raise AssertionError("SKILL.md must start with YAML frontmatter")

    try:
        closing_index = lines.index("---", 1)
    except ValueError as error:
        raise AssertionError("SKILL.md frontmatter must be closed") from error

    fields: dict[str, str] = {}
    for line in lines[1:closing_index]:
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        key, separator, value = line.partition(":")
        if not separator or not key.strip() or not value.strip():
            raise AssertionError(f"unsupported frontmatter line: {line!r}")
        fields[key.strip()] = value.strip().strip("\"'")
    return fields, "\n".join(lines[closing_index + 1 :])


def parse_scalar_yaml(document: str) -> dict[str, str | bool]:
    fields: dict[str, str | bool] = {}
    section = ""
    for line in document.splitlines():
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue

        indent = len(line) - len(line.lstrip(" "))
        key, separator, raw_value = stripped.partition(":")
        if not separator:
            raise AssertionError(f"unsupported openai.yaml line: {line!r}")
        if not raw_value.strip():
            if indent != 0:
                raise AssertionError(f"unsupported nested section: {line!r}")
            section = key
            fields[section] = ""
            continue

        path = f"{section}.{key}" if indent else key
        if indent == 0:
            section = ""
        value = raw_value.strip().strip("\"'")
        if value == "true":
            fields[path] = True
        elif value == "false":
            fields[path] = False
        else:
            fields[path] = value
    return fields


def read_required(path: Path) -> str:
    if not path.is_file():
        raise AssertionError(f"missing {path.relative_to(ROOT)}")
    return path.read_text(encoding="utf-8")


def parse_contract(body: str) -> element_tree.Element:
    start_token = "<next-skill-contract>"
    end_token = "</next-skill-contract>"
    start = body.find(start_token)
    end = body.find(end_token)
    if start < 0 or end < 0:
        raise AssertionError("SKILL.md must contain a next-skill-contract block")
    xml = body[start : end + len(end_token)]
    return element_tree.fromstring(xml)


class TestSkillContract(unittest.TestCase):
    def test_skill_artifacts_exist_when_v02_contract_is_added(self) -> None:
        self.assertTrue(SKILL_PATH.is_file(), f"missing {SKILL_PATH.relative_to(ROOT)}")
        self.assertTrue(OPENAI_PATH.is_file(), f"missing {OPENAI_PATH.relative_to(ROOT)}")

    def test_frontmatter_routes_only_explicit_next_without_execution(self) -> None:
        fields, _body = parse_frontmatter(read_required(SKILL_PATH))

        self.assertEqual(fields.get("name"), "next")
        description = fields.get("description", "").lower()
        self.assertIn("$next", description)
        self.assertIn("never execute", description)

    def test_openai_metadata_disables_implicit_invocation_and_dependencies(self) -> None:
        fields = parse_scalar_yaml(read_required(OPENAI_PATH))

        self.assertTrue(fields.get("interface.display_name"))
        self.assertTrue(fields.get("interface.short_description"))
        self.assertTrue(fields.get("interface.default_prompt"))
        self.assertIs(fields.get("policy.allow_implicit_invocation"), False)
        self.assertFalse(any(path == "dependencies" or path.startswith("dependencies.") for path in fields))

    def test_manifest_declares_v02_skill_only_package_surface(self) -> None:
        manifest = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))

        with self.subTest(field="version"):
            self.assertEqual(manifest.get("version"), "0.2.0")
        with self.subTest(field="skills"):
            self.assertEqual(manifest.get("skills"), "./skills/")
        for forbidden_field in ("hooks", "mcpServers", "apps"):
            self.assertNotIn(forbidden_field, manifest)

    def test_body_structurally_limits_inputs_and_side_effects(self) -> None:
        _fields, body = parse_frontmatter(read_required(SKILL_PATH))
        contract = parse_contract(body)

        self.assertEqual(contract.findtext("context"), "current-context-only")
        self.assertEqual(contract.findtext("tools"), "none")
        self.assertEqual(contract.findtext("sources/transcript"), "forbidden")
        self.assertEqual(contract.findtext("sources/history"), "forbidden")
        self.assertEqual(contract.findtext("sources/model-api"), "forbidden")
        self.assertEqual(contract.findtext("execution"), "forbidden")

    def test_body_structurally_defines_bounded_recommendations(self) -> None:
        _fields, body = parse_frontmatter(read_required(SKILL_PATH))
        contract = parse_contract(body)

        self.assertEqual(contract.findtext("recommendations/default"), "1")
        self.assertEqual(contract.findtext("recommendations/maximum"), "3")
        self.assertEqual(contract.findtext("recommendations/forks"), "real-only")
        self.assertEqual(contract.findtext("language"), "same-as-user")

    def test_body_structurally_rejects_implicit_surfaces(self) -> None:
        _fields, body = parse_frontmatter(read_required(SKILL_PATH))
        contract = parse_contract(body)

        self.assertEqual(contract.findtext("invocation"), "explicit-$next-only")
        self.assertEqual(contract.findtext("custom-slash-command"), "none")
        self.assertEqual(contract.findtext("automatic-footer"), "forbidden")


if __name__ == "__main__":
    unittest.main()
