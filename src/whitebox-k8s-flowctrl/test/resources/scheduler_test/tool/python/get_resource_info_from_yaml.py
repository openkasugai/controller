#!/usr/bin/env python3
# Copyright 2024 NTT Corporation , FUJITSU LIMITED
import argparse
import sys
import yaml

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--yamlfile", help="YAML File")
    args = parser.parse_args()

    exitCode = 0
    if args.yamlfile:
        with open(args.yamlfile) as f:
            exitCode = getResourceInfo(f)
    else:
        exitCode = getResourceInfo(sys.stdin)

    return exitCode

def getResourceInfo(f):
    try:
        yml = yaml.safe_load(f)
        if 'kind' in yml and 'metadata' in yml:
            kind = yml.get('kind')
            metadata = yml.get('metadata')
            name = metadata.get('name')
            namespace = metadata.get('namespace')
            print(kind, name, namespace)
            return 0
        else:
            print("kind/metadata does not exist.")
            return 1
    except Exception as exc:
        print("failed parsing yaml.", exc, file=sys.stderr)
        return False

if __name__ == "__main__":
    sys.exit(main())
