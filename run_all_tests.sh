GOPATH="$PWD"/"`dirname $0`"

echo "Setting GOPATH to $GOPATH"
export GOPATH

echo "Running Go tests"

go test -v echodiscovery/... -timeout 10s

rc=$?
if [ $rc -ne 0 ]; then
  echo "Test failures detected!"
else
  echo "Tests succeeded!"
fi
exit $rc
