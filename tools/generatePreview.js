const fs = require('fs');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;
const generate = require('@babel/generator').default;
const t = require('@babel/types');

const sourceCode = fs.readFileSync(process.argv[2], 'utf8');

const ast = parser.parse(sourceCode, {
    sourceType: 'module',
    plugins: ['typescript', 'jsx']
});

traverse(ast, {
    ReturnStatement(path) {
        const arg = path.node.argument;

        if (t.isJSXElement(arg) && !(arg.openingElement.attributes.length > 0 && arg.openingElement.attributes[0].value.value === "border border-red-500")) {
            const wrapperDiv = t.jSXElement(
                t.jSXOpeningElement(
                    t.jSXIdentifier('div'),
                    [
                        t.jSXAttribute(
                            t.jSXIdentifier('className'),
                            t.stringLiteral('border border-red-500')
                        )
                    ],
                    false
                ),
                t.jSXClosingElement(t.jSXIdentifier('div')),
                [arg],
                false
            );

            path.node.argument = wrapperDiv;
        }
    }
});

const { code } = generate(ast, {}, sourceCode);

fs.writeFileSync(process.argv[2], code);
