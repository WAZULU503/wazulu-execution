# remove the duplicate hash() function from verify.go automatically
awk '
BEGIN{skip=0}
{
  if(skip==0 && $0 ~ /^func hash\(/){skip=1; brace=0}
  if(skip==1){
      brace+=gsub(/{/,"{")
      brace-=gsub(/}/,"}")
      if(brace<=0){skip=0; next}
      next
  }
  print
}' log/verify.go > log/verify.go.tmp && mv log/verify.go.tmp log/verify.go

# rebuild
go build -o wz ./cmd/wz

# run program
./wz exec
./wz verify
