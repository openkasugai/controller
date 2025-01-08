#!/usr/bin/env python3
# Copyright 2024 NTT Corporation , FUJITSU LIMITED
import argparse
import sys
import yaml


def main():
    parser = argparse.ArgumentParser(prog='update_status.sh')
    subparsers = parser.add_subparsers(help="sub-command", required=True)

    parserUpdate = subparsers.add_parser('update')
    parserUpdate.add_argument("kind", help="Resource Kind")
    parserUpdate.add_argument("name", help="Resource Name")
    parserUpdate.add_argument("-n", "--namespace", help="Resource Namespace (default: %(default)s)", default="default")
    parserUpdate.add_argument("-f", "--file", help="YAML File", default="")

    parserApply = subparsers.add_parser('apply')
    parserApply.add_argument("-f", "--file", help="YAML File", required=True)

    try:
        args = parser.parse_args()

        if hasattr(args, 'kind'):
            print(args.kind, args.name, args.namespace, args.file)
        else:
            print(args.file)
        return 0
    except Exception as exc:
        print(exc)
        return 1

if __name__ == "__main__":
    sys.exit(main())
