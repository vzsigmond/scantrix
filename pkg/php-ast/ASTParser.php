<?php
/**
 * @file ASTParser.php
 *
 * Usage: php ASTParser.php /path/to/phpfile.php
 */

class ASTParser {
    private $version;

    /**
     * Constructor.
     *
     * @param int $version The AST version to use (e.g. 50)
     */
    public function __construct($version = 50) {
        $this->version = $version;
    }

    /**
     * Parses the provided PHP code and returns the AST.
     *
     * @param string $code The PHP code to parse.
     * @return ast\Node The root AST node.
     */
    public function parse($code) {
        return ast\parse_code($code, $this->version);
    }

    /**
     * Recursively converts an AST node into an associative array.
     *
     * @param mixed $node An AST node or scalar.
     * @return mixed The node converted to an array, or the scalar.
     */
    public function astToArray($node) {
        if (!($node instanceof ast\Node)) {
            return $node; // Return scalars as is (e.g. string, int, null).
        }
        $result = [
            'kind'     => ast\get_kind_name($node->kind),
            'flags'    => $node->flags,
            'lineno'   => $node->lineno,
            'children' => [],
        ];
        foreach ($node->children as $key => $child) {
            $result['children'][$key] = $this->astToArray($child);
        }
        return $result;
    }

    /**
     * Parses the PHP code and returns the AST as a JSON-formatted string.
     *
     * @param string $code The PHP code to parse.
     * @return string The JSON representation of the AST.
     */
    public function parseToJson($code) {
        $ast = $this->parse($code);
        $arrayAst = $this->astToArray($ast);
        return json_encode($arrayAst, JSON_PRETTY_PRINT);
    }
}

// Main Execution

if ($argc < 2) {
    // No file path provided; output usage to STDERR and exit.
    fwrite(STDERR, "Usage: php ast_script.php /path/to/phpfile.php\n");
    exit(1);
}

$filePath = $argv[1];
if (!file_exists($filePath)) {
    fwrite(STDERR, "Error: File '$filePath' does not exist.\n");
    exit(1);
}

$code = file_get_contents($filePath);

$parser = new ASTParser(110);
echo $parser->parseToJson($code);
