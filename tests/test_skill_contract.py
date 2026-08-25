import unittest
import xml.etree.ElementTree as element_tree
from pathlib import Path
from typing import Final


ROOT: Final = Path(__file__).resolve().parents[1]
SKILL_PATH: Final = ROOT / "skills" / "next" / "SKILL.md"
OPENAI_PATH: Final = ROOT / "skills" / "next" / "agents" / "openai.yaml"


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
    def test_standalone_skill_contains_only_skill_markdown(self) -> None:
        self.assertTrue(SKILL_PATH.is_file(), f"missing {SKILL_PATH.relative_to(ROOT)}")
        self.assertFalse(OPENAI_PATH.exists(), f"remove {OPENAI_PATH.relative_to(ROOT)}")
        skill_files = tuple(
            path.relative_to(SKILL_PATH.parent).as_posix()
            for path in sorted(SKILL_PATH.parent.rglob("*"))
            if path.is_file()
        )
        self.assertEqual(skill_files, ("SKILL.md",))

    def test_frontmatter_routes_only_explicit_next(self) -> None:
        fields, _body = parse_frontmatter(read_required(SKILL_PATH))

        self.assertEqual(fields.get("name"), "next")
        description = fields.get("description", "").lower()
        self.assertIn("$next", description)
        self.assertIn("literal", description)
        self.assertIn("only", description)
        self.assertIn("do not use", description)
        self.assertIn("what should i do next", description)
        self.assertIn("current conversation", description)
        self.assertNotIn("never execute", description)

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

    def test_body_defines_prompt_only_output_without_placeholders(self) -> None:
        _fields, body = parse_frontmatter(read_required(SKILL_PATH))
        contract = parse_contract(body)

        self.assertEqual(contract.findtext("output/format"), "fenced-text")
        self.assertEqual(contract.findtext("output/commentary"), "forbidden")
        self.assertEqual(contract.findtext("output/placeholders"), "forbidden")


if __name__ == "__main__":
    unittest.main()
