#!/usr/bin/env -S python3 -S
from typing import Any
import argparse
import pathlib
import sys
import json
import subprocess

_LOCATION_TABLE_HEADERS = (
    "type",
    "scene",
    "default",
    "addresses",
    "vanilla",
    "categories",
)


def main(args: argparse.Namespace) -> int:
    sys.path.append(args.zootr)
    from version import __version__
    from Item import ItemInfo
    import LocationList

    output = pathlib.Path(args.output)
    zootr = pathlib.Path(args.zootr)

    if output.exists():
        print(f"Cleaning {output}")
        subprocess.run(["rm", "-r", str(output)])

    output.mkdir(exist_ok=True)
    (output / "data").mkdir(exist_ok=True)
    (output / "logic" / "glitchless").mkdir(parents=True, exist_ok=True)
    (output / "logic" / "glitched").mkdir(exist_ok=True)

    print(f"Dumping from OOT-Randomizer v{__version__}")
    dump_to_file("items", output / "data" / "items.json",
                 list(ItemInfo.items.values()))
    dump_to_file(
        "locations",
        output / "data" / "locations.json",
        rearrange_location_table(LocationList.location_table),
    )
    copy_logic_dir(zootr / "data" / "World", output / "logic" / "glitchless")
    copy_logic_dir(zootr / "data" / "Glitched World",
                   output / "logic" / "glitched")
    copy_file(zootr / "data" / "LogicHelpers.json",
              output / "logic" / "helpers.json")
    return 0


def cli() -> argparse.ArgumentParser:
    cli = argparse.ArgumentParser(
        "zootr-dump", description="Dumps data from OOT-Randomizer")
    cli.add_argument(
        "-Z",
        "--zootr",
        help="Path to zootr source",
        required=True,
        type=str,
        dest="zootr",
    )

    cli.add_argument(
        "-O",
        "--output",
        help="Path to dump files",
        required=True,
        type=str,
        dest="output",
    )

    return cli


def dump_to_file(kind: str, path: pathlib.Path, obj: Any):
    print(f"Dumping {kind} to {path}")
    with open(path, mode="w") as fh:
        json.dump(obj, fh, cls=ZootrJsonEncoder, indent=2)


def _rename_file(filename: str) -> str:
    return filename.lower().replace(' ', '-')


def copy_logic_dir(src: pathlib.Path, dest: pathlib.Path) -> None:
    for f in src.iterdir():
        copy_file(f, dest / _rename_file(f.parts[-1]))


def copy_file(src: pathlib.Path, dest: pathlib.Path) -> None:
    print(f"Copying {src} to {dest}")
    dest.touch(exist_ok=True)
    with open(src, mode="r") as s, open(dest, mode="w") as d:
        d.write(s.read())


def rearrange_location_table(tbl) -> Any:
    return [{
        "name": name,
        **dict(
            filter(
                lambda p: p[0] not in {"scene", "default", "addresses"},
                zip(_LOCATION_TABLE_HEADERS, data)
            ))}
            for name, data in tbl.items()
            ]


class ZootrJsonEncoder(json.JSONEncoder):
    def default(self, o: Any) -> Any:
        from Item import ItemInfo

        match o:
            case ItemInfo():
                return self._dump_item_info(o)
            case _:
                return super().default(o)

    def _dump_item_info(self, o: Any) -> Any:
        special = {k: v for k, v in o.special.items()}

        if (alias := special.pop("alias", None)) != None:
            special["alias"] = {"name": alias[0], "value": alias[1]}

        if special.get("progressive") == float("Inf"):
            special["progressive"] = -1

        if special.get("bottle") == float("Inf"):
            special["bottle"] = -1

        return {
            "name": o.name,
            "advancement": o.advancement,
            "priority": o.priority,
            "type": o.type,
            "special": special,
        }


if __name__ == "__main__":
    parser = cli()
    if (code := main(parser.parse_args())) != 0:
        parser.print_help()
    sys.exit(code)
