#!/bin/zsh
#encoding=utf-8

function run_server () {
    # param: sleep_time
    nohup python3 -u "/Users/tomjack/Desktop/code/Python/SHU_report_public/shu_f/server.py" >> serverlog.log &
    # sleep $1
}

function initServer () {
    # param: sleep_time
    init_res=$(curl -s 127.0.0.1:8989/init)
    printf '%s\n' "$init_res"
    sleep $1
    check_res=$(curl -s 127.0.0.1:8989/test)
    if [ "${check_res}" = '"about:blank"' ]; then
        echo "Browser launch successfully."
    fi
}

function main() {
    run_server 
    # initServer 1
}

main