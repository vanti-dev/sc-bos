#!/bin/sh
#
# PROVIDE: areacontroller
# REQUIRE: networking TcSystemService
# KEYWORD:

. /etc/rc.subr

name="areacontroller"
rcvar="areacontroller_enable"

ops_chdir="/usr/local/opt/areacontroller"

load_rc_config $name
: "${areacontroller_enable:="no"}"
: "${areacontroller_user:="areacontroller"}"
: "${areacontroller_group:="areacontroller"}"
: "${areacontroller_flags:="-data-dir /usr/local/opt/areacontroller/data -static-dir /usr/local/opt/areacontroller/static"}"
: "${areacontroller_log:="/var/log/areacontroller.log"}"

areacontroller_cmd="/usr/local/opt/areacontroller/bin/areacontroller ${areacontroller_flags}"

pidfile="/var/run/${name}.pid"
# command="/usr/sbin/daemon"
command="/usr/sbin/daemon"
areacontroller_daemon_args="-c -f -o ${areacontroller_log} -P ${pidfile} -r -H"

start_precmd="areacontroller_precmd"
areacontroller_precmd()
{

        if [ ! -e "${areacontroller_log}" ]; then
            install -g ${areacontroller_group} -o ${areacontroller_user} -- /dev/null "${areacontroller_log}";
        fi
        install -o "${areacontroller_user}" /dev/null "${pidfile}"
        rc_flags="${areacontroller_daemon_args} ${areacontroller_cmd} ${rc_flags}"
}

run_rc_command "$1"