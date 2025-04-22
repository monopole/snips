#!/bin/bash

# quit on any error
set -e

# Get a "Personal Access Token" (aka bearer token) from your jira host
# at https://issues.acmecorp.com/secure/ViewProfile.jspa
if [ -z "$JIRA_API_TOKEN" ]; then
  echo "Set JIRA_API_TOKEN"
fi

# The following doesn't work:
password="elided"
basic=$(echo "jregan@acmecorp.com:$password" | base64 --wrap=0 )
# Don't use
#   Authentication: Basic {encoded email : bearer},
#   See
#   https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html
#
# Instead use
#    Authentication: Bearer ${rawRat}

# --insecure \
# --cacert "~/.acmecorpCert.pem" \
# --silent --show-error \
# -H "accept: application/json" \

host="https://issues.teslamotors.com"
dataForPost=/tmp/post.json
parsedData=/tmp/parsedFromPost.json
searchEndpoint="rest/api/2/search"

output=/mnt/c/Users/jregan/Downloads/k.json
rm -f $output

echo "basic = $(echo $basic | base64 -d)"


function authHeader {
  echo  "Authorization: Bearer $JIRA_API_TOKEN"
#  echo  "Authorization: Basic $basic"
}

echo "$(authHeader)"

function doGet {
  curl \
  --insecure \
  --silent \
  --show-error \
  --request GET \
  --header "Accept: application/json" \
  --header "$(authHeader)" \
  --output $output \
  "$host/$1"
}

#  --silent \
#  --show-error \

function doPost {
  echo "Check data in ${2}:"
  # the following assures clean data
  jq . $2 >$parsedData
  echo "No error in jq use; posting data to $1"
  curl \
  --insecure \
  --request POST \
  --data @$parsedData \
  --header "Content-Type: application/json" \
  --header "Accept: application/json" \
  --header "$(authHeader)" \
  --output $output \
  "$host/$1"
  echo "curl result = $?"
}



function doGet1 {
  doGet "rest/api/2/issue/PLM-19711"
}
function doPost1 {
  cat <<'EOF' >$dataForPost
{
  "jql": "order by created DESC",
  "maxResults": 2,
  "fields": ["status"],
  "expand": ["changelog", "renderedFields", "names", "schema", "transitions", "operations", "editmeta"]
}
EOF
    doPost $searchEndpoint $dataForPost
}

function doGet2 {
  doGet "${searchEndpoint}?jql=assignee=jregan+order+by+duedate&fields=id,key&maxResults=5"
}

function doPost2 {
  cat <<'EOF' >$dataForPost
{
  "jql": "issuefunction in commented ('by jregan after 2023/03/01') and issuefunction in commented ('by jregan before 2023/04/01')",
  "maxResults": 2,
  "fields": ["id", "key"],
  "expand": ["renderedFields", "names", "schema"]
}
EOF
    doPost ${searchEndpoint} $dataForPost
}

function doPost3 {
  cat <<'EOF' >$dataForPost
{
  "jql": "creator = jregan and created >= '2023/06/08' and created <= '2023/06/08'",
  "maxResults": 2,
  "fields": ["id", "key"],
  "expand": ["renderedFields", "names", "schema"]
}
EOF
    doPost $searchEndpoint $dataForPost
}

function doPost3 {
  #  "fields": ["id", "key"],
 cat <<'EOF' >$dataForPost
{
  "jql": "creator = jregan and created >= '2023/06/08' and created < '2023/06/09'",
  "maxResults": 2,
  "fields": [
    "id",
    "key",
    "summary",
    "resolution",
    "labels",
    "assignee",
    "reporter",
    "project",
    "description",
    "creator",
    "updated"
  ]
}
EOF
    doPost $searchEndpoint $dataForPost
}
#  "expand": ["renderedFields", "names", "schema" ]

function doPost4 {
  cat <<'EOF' >$dataForPost
{
  "jql": "creator = jregan and created during ('2023/06/07','2023/06/09')",
  "maxResults": 2,
  "fields": ["id", "key", "summary", "resolution", "labels", "assignee", "reporter", "project", "description","creator","updated"],
  "expand": ["renderedFields", "names", "schema"]
}
EOF
    doPost $searchEndpoint $dataForPost
}

# These all work
#doGet1
#doGet2
#doPost1
#doPost2 # works
doPost3 # works
#doPost4 # does not work, DURING isn't recognized

echo " "
echo " "
echo "Result $?"
echo "Input: ======"
jq . $dataForPost
echo "Output: ======"
jq . $output
echo "======"
echo "jq . $output | more"
