# Website
This is a small custom lexer i made in order to generate my own static website
The code is not very clean and its uncommented but it works
I dont recommend using it for anything else than my own website
But feel free to look around i guess

# How it works

The lexer first gets all the htm files in the directory and then it parses them
it does this by reading through each character and sends it through a chain of nodes that are working within eachother.
basically the first node is the raw node, once it reads a character it checks with each viable node if it can be considered a token, if it can the subnode will be called and blocks the flow from continuing and scans further using the same node and does the same.