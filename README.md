# go-simple-expression-eval
Simple expression evaluator in golang for learning lexers, parsers, ast's and interpreters.   

This project doesn't support unary operators, so expression "-1+2" will report us "Missing operand".
For the unary operator support a modified version of Shunting-yard algorithm is needed ([See this project](https://github.com/MacTee/Shunting-Yard-Algorithm/blob/master/ShuntingYard/InfixToPostfixConverter.cs))

# Running
**go build**  
**./go-simple-expression-eval "1+2+3*5-(5-8)"**

// quote your expression as it might contain shell operator:  

# Program flow
1. Tokenize (Lexer) 
2. Convert **infix form** to **postfix form**  using Shunting-yard algorithm (Parser)
3. Construct **Abstract syntax tree** - AST (Parser)
4. Evaluate AST (Interpreter)

# Lexer
_lexer.go : Lex()_

Purpose of lexer is to turn input into sequence of tokens (lex items) and determine their type while thrashing whitespaces and reporting unrecognized characters.

Lexer on its own doen't know grammar of parsed language so it does't report syntactic errors or semantic errors.

I heavily recommend watching [Youtube - Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE) as I use these principles.

Given well-formed expression input: "1+2+3", lexer correctly produces sequence "1", "+", "2", "+", "3".

Given not well-formed expression, such as "1)+3(((" lexer **doesn't report any errors**, but produces sequence "1", ")", "+", "3", "(", "(", "(". Errors will be reported in parsing phase not in lexing phase.

Gived input with unrecognized characters "*12+死+3" lexer produces error:

**Lexing error at 4: Invalid symbol: '死'**

Our lexer also determines type and position of lexed items, which will come handy later.  
Given expression "1+1*(5-5)" lexer produces:

Type: INUMBER, Val: "1", Pos: 1  
Type: IADD, Val: "+", Pos: 2  
Type: INUMBER, Val: "1", Pos: 3  
Type: IMUL, Val: "*", Pos: 4  
Type: ILPAR, Val: "(", Pos: 5  
Type: INUMBER, Val: "5", Pos: 6  
Type: ISUB, Val: "-", Pos: 7  
Type: INUMBER, Val: "5", Pos: 8  
Type: IRPAR, Val: ")", Pos: 9  
Type: EOF, Val: "", Pos: 9  



# Form conversion
_parser.go : toPostfix()_

The classical format of expressions "&lt;operand&gt; &lt;operator&gt; &lt;operand&gt;" is called **infix form**. This form is somewhat troublesome when constructing AST so first we will convert **Infox form** to **Postfix form** using Shunting yard algorithm. Postfix form uses rather weird syntax "&lt;operand&gt; &lt;operand&gt; &lt;operator&gt;". The advantage is that we get rid of the need for parentheses (those are required in infix form).

This algorithm goes trough out lexer tokens and using **stack data structure** and **list** converts the form. Stack is used as stacked operator store and list for postfix output.

The algorithm is following:
1. If current token is number (INUMBER) then put it to output
2. If current token is left parenthesis then Push it to stack
3. If current token is right parenthesis then Pop stack until we hit left parenthesis and Pop it out
4. If current token is operator +,-,*,/ Pop the stack until we find operator with higher or equal precedence (see below) and Push current operator
5. When we go trough all tokens, Pop all items from stack to output

**Operator precedence:**  
Add and Subtract = 1  
Multiply and Divide = 2

For more about Postfix form see:  
[Wikipedia - Reverse Polish notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation)  
[Expression binary trees and forms](http://www.cim.mcgill.ca/~langer/250/19-binarytrees-slides.pdf)

# Constructing Abstract syntax tree - AST
_parser.go : constructAst()_

AST for expressions is a simple binary tree. Once we have our expression in postfix form, creating AST is trivial. For construction we will go trough our postfix formed lex items and using stack we build a tree.

In the tree we distinguish types of nodes we will use later in interpretation. Nodes containing operator (+-*/) hold no value, just left ,right operand children. Nodes containing numbers have no children (leaves).

Algorithm for AST from postfix form is following:  
Go trough all postfix items  
1. If current item is number, create node with this value and Push it to stack.
2. If its operator, create node for it, then Pop the stack and put it to right child, Pop stack again and put it to left child (order is important). Then push new operator node to stack.  

When there are no more postfix items, Pop the stack last time, this is your AST root node.

# Interpreting (Evaluating AST)
_interpreter.go : Interpret(), postOrderTraversal()_


Once we have AST the interpreting is childs play. All we need to do is do a recursive **Post order traversal** of our AST binary tree and execute corresponding operation.

Post order traversal interpretation algorithm:
1. Start at root
2. If current node is leaf node return its value
3. Recurse to left child and store return value
4. Recurse to right child and store return value
5. Perform arithmetic operation on stored left and right values
6. Return result of operation

Result of traversal is calculated value.

# Sources
[Youtube - Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)  
[Lexer theory](http://www.cse.chalmers.se/edu/year/2010/course/TIN321/lectures/proglang-04.html)  
[Binary and expression trees](http://www.cim.mcgill.ca/~langer/250/19-binarytrees-slides.pdf)  
[Minimal expression evaluator in Go](https://rosettacode.org/wiki/Arithmetic_Evaluator/Go)  
[More complex lexer for template engine](https://golang.org/src/text/template/parse/lex.go)  
[AST creation](https://softwareengineering.stackexchange.com/questions/254074/how-exactly-is-an-abstract-syntax-tree-created)  





