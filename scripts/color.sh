#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'
YELLOW='\033[0;33m'
NC2='\033[4m'
PURPLE='\033[0;35m'


function print_blue() {
  printf "${BLUE}%s${NC}\n" "$1"
}

function print_green() {
  printf "${GREEN}%s${NC}\n" "$1"
}

function print_red() {
  printf "${RED}%s${NC}\n" "$1" 
}

function print_yellow() {
  printf "${YELLOW}%s${NC}\n" "$1"
}

function print_yellow2() {
  printf "${YELLOW}%s${NC2}\n" "$1"
}

function print_purple() {
  printf "${PURPLE}%s${NC}\n" "$1"
}


<<comment
字背景颜色范围: 40--49                   字颜色: 30--39  
            40: 黑                          30: 黑  
        41:红                          31: 红  
        42:绿                          32: 绿  
        43:黄                          33: 黄  
        44:蓝                          34: 蓝  
        45:紫                          35: 紫  
        46:深绿                        36: 深绿  
        47:白色                        37: 白色


输出特效格式控制：  
  
\033[0m  关闭所有属性    
\033[1m   设置高亮度    
\03[4m   下划线    
\033[5m   闪烁    
\033[7m   反显    
\033[8m   消隐    
\033[30m   --   \033[37m   设置前景色    
\033[40m   --   \033[47m   设置背景色  
comment

<<comment
print_blue "Usage:  "
echo "why you are here?"
print_red  "red:    "
print_green "green:  "
print_yellow "yellow:this is yelloe "
print_yellow2 "you can see it"
  echo "  playground.sh <mode>"
  echo "    <OPT> - one of 'up', 'down', 'restart'"
  echo "      - 'up' - bring up the bitxhub network"
  echo "      - 'down' - clear the bitxhub network"
  echo "      - 'restart' - restart the bitxhub network"
  echo "  playground.sh -h (print this message)"
comment

