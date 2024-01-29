#Runs all tests in the names packages named in ./testpackages
file="./test_packages"

# Check if the file exists
if [ ! -f "$file" ]; then
    echo "File not found: $file"
    exit 1
fi

mapfile -t dirs < "$file"

for dir in "${dirs[@]}"; do
    echo "Testing in: $dir" 
    go test -race $dir
done

echo "Tests completed"
