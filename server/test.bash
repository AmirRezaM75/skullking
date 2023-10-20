#!/bin/bash
for i in {1..500}
do
   curl --location --request POST 'localhost:3000/games' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDAiLCJleHAiOjE2OTYwMTU3NzEsImp0aSI6IjY1MTZmYjZiYzhiZjk2MTliZDUwMTFkNiJ9.QgtjvQH9bpZGhWjAghwnqYBXkq1dDe9KnlN9X3omDSQ'
done
