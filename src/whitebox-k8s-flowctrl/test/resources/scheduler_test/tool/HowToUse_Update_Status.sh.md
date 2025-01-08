## Prerequisites
status を設定する対象のカスタムリソースが登録済みであること。

## composition
```
update_status.sh                        : Status更新スクリプト
python/  
  ├ parse_args.py                       : 引数を解析するスクリプト
  ├ get_status_yaml.py                  : YAMLからStatus部分のyamlを取り出すスクリプト
  ├ yaml2json.py                        : YAMLをJSONに変換するスクリプト
  └ get_resource_info_from_yaml.py      : YAMLからカスタムリソースのkind, name, namespace を取得するスクリプト
```

 ## Instructions
 以下の1,2のいずれかの方法で update_status.sh を実行してください。
 ```
1. update_status.sh update kind name [-n NAMESPACE] [-f FILE]  
    - kind          :カスタムリソースのkind(複数形の名称は不可)  
    - name          :カスタムリソース名
    - NAMESPACE     :名前空間(省略字は"default"を使用)
    - FILE          :statusを記載したyamlファイルのパス(省略時は標準入力)  
                     status以外の項目は無視
2. update_status.sh apply -f FILE
    - FILE  :statusファイルを記載したyamlファイルのパス  
             ただし、kind, metadata.name, metadata.namespace の記載も必要
```
## Notice
コピーして使う場合は、"構成"に記載の全ファイルをディレクトリ構成を変えずにコピーすること。