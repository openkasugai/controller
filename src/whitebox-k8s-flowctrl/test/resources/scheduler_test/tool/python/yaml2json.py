#!/usr/bin/env python3
# Copyright 2024 NTT Corporation , FUJITSU LIMITED
import yaml
import json
import sys

def main():
    return yamlToJson(sys.stdin)

def yamlToJson(yamlStr):
    try:
        obj = yaml.load(yamlStr, Loader=yaml.SafeLoader)
        json.dump(obj, sys.stdout, indent=2)
        return True
    except Exception as exc:
        print("failed to convert yaml to json.", exc, file=sys.stderr)
        return False

if __name__ == "__main__":
    if main(): 
        sys.exit(0)
    else:
        sys.exit(1)
