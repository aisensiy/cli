#!/usr/bin/env bash
user_name="user`date +%s`"

watchError(){
  while read data
  do
    if [[ $data == *"Error"* ]]
    then
      echo "Fail"
      break
    fi
  done
}

watchWhoamI(){
  read data
  if [[ $data != *"${user_name}"* ]]
  then
    echo "Fail";
  fi
}

stack="javajersey-test"
watchStack() {
  found=false
  while read line
  do
    if [[ $line == *"${stack}"* ]]
    then
      found=true
    fi
  done

  if [[ ${found} != true ]]
  then
    echo "Fail";
  fi
}

watchApp(){
  found=false
  while read line
  do
    if [[ $line == *"${appName}"* ]]
    then
      found=true
    fi
  done

  if [[ ${found} != true ]]
  then
    echo "Fail";
  fi
}

watchDomain() {
  found=false
  while read line
  do
    if [[ $line == *"Domains ["[1-9]*"]"* ]]
    then
      found=true
      break
    fi
  done

  if [[ ${found} != true ]]
  then
    echo "Fail";
  fi
}

routeId="origin"
watchNewRoute() {
  while read data
  do
    if [[ $data == *"Error"* ]]
    then
      echo "Fail"
      return
    elif [[ $data==*"http"* ]]
    then
      routeId=${data##*/}
      break
    fi
  done

  if [[ ${routeId} == "" ]]
  then
    echo "Fail"
    return
  fi
  echo ${routeId}
}

findInLines() {
  found=false
  while read line
  do
    if [[ $line == *"$1"* ]]
    then
      found=true
      break
    fi
  done

  if [[ ${found} != true ]]
  then
    echo "Fail";
  fi
}

watchRoutes() {
  found=false
  while read line
  do
    if [[ $line == *"Routes: ["[1-9]*"]"* ]]
    then
      found=true
      break
    fi
  done

  if [[ ${found} != true ]]
  then
    echo "Fail";
  fi
}

watchKeys() {
  findInLines "id_rsa"
}

sshKeyId=""
watchNewKey() {
  while read line
  do
    if [[ $line == *"Error"* ]]
    then
      echo "Fail";
      return
    elif [[ $line == *"http"* ]]
    then
      sshKeyId=${line##*/}
      break
    fi
  done

  if [[ ${sshKeyId} == "" ]]
  then
    echo "Fail"
  fi
}

report(){
  read data
  if [[ $data != "Fail" ]]
  then
    printf '\e[1;32m%-6s\e[m \n' "CMD $1 : Ok"
  else
    printf '\e[1;31m%-6s\e[m' "CMD $1 : Fail"
  fi
}

cde="./cde"
appName="app`date +%s`"
stackName="stack`date +%s`"
domainName="domain.`date +%s`.com"
routePath="path/`date +%s`"

${cde} register http://192.168.50.4:31088 --email $user_name@tw.com --password admin 2>&1\
| watchError | report "register"
${cde} login http://192.168.50.4:31088 --email $user_name@tw.com --password admin\
 2>&1 | watchError | report "login"
${cde} whoami 2>&1 | watchWhoamI | report "whoami"

${cde} stacks:create ${stackName} 2>&1 | watchError | report "stacks:create"
${cde} stacks:list 2>&1 | watchStack | report "stacks:list"
${cde} stacks:remove  2>&1 | watchError | report "stacks:remove"

sshKeyId=`${cde} keys:add ~/.ssh/id_rsa.pub 2>&1 | watchNewKey `; echo $sshKeyId| report "keys:add"
if [[ ${sshKeyId} == "Fail" ]]; then
  exit 1
fi
${cde} keys:list 2>&1 | watchKeys | report "keys:list"
${cde} keys:remove ${sshKeyId} 2>&1 | watchError | report "keys:remove"

# prepare stack data
stackName="stackjava"
${cde} stacks:create ${stackName} > /dev/null 2>&1
git remote remove cde
${cde} apps:create ${appName} ${stackName} 2>&1 | watchError | report "apps:create"
${cde} apps:list 2>&1 | watchApp | report "apps:list"
${cde} apps:info -a ${appName} 2>&1 | watchApp | report "apps:info"

${cde} domains:add ${domainName} 2>&1 | watchError | report "domains:add"
${cde} domains:list 2>&1 | watchDomain | report "domains:list"
${cde} domains:remove ${domainName} 2>&1 | watchError | report "domains:remove"

# prepare domain data
domainName="domain.template.com"
${cde} domains:add ${domainName} > /dev/null 2>&1
routeId=`${cde} routes:create ${domainName} ${routePath} 2>&1 | watchNewRoute `; echo $routeId | report "routes:create"
${cde} routes:list 2>&1 | watchRoutes | report "routes:list"
if [[ ${routeId} == "Fail" ]]; then
  exit 1
fi
${cde} routes:bind ${routeId} ${appName} | watchError | report "routes:bind"
${cde} routes:unbind ${routeId} ${appName} | watchError | report "routes:unbind"

# destroy resources
# user unregister
