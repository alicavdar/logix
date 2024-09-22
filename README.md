# Logix: A simple rule-based condition language

Logix is a lightweight language for defining and checking conditions using custom fields and values. It supports different comparison operators, string matching, array checks, and logical groupings, so you can easily build complex rule sets.

Logix is indentation sensitive, similar to Python, meaning indentation is used to group conditions. It also supports comments, which can be added with the # symbol.

## Supported Operators

- Basic comparisons: `eq`, `neq`, `gt`, `lt`, `gte`, and `lte`.
- String operations: `contains`, `startsWith`, `endsWith`.
- Range checks: Use `between` to see if a value is in a certain range.
- Logical operators: Combine conditions using `and` and `or`.
- Array checks: Use `in` to check if a value exists in a list.
- Negation: Use `not` to negate `in`, `contains`, `between`, `startsWith`, and `endsWith` operators.

Logix supports deeply nested fields like `products[0].info.title` and valid boolean and null values such as `true`, `false`, and `nil` in conditions.

Here's an example of how Logix syntax looks:

```
group and
    price gt 100
    status eq "active"
    # This is a comment
    group or
        category in ["electronics", "furniture"]
        stock between 50 and 100
    title not contains "deleted"
```

## Usage

Here’s a context example:

```go
context := map[string]interface{}{
    "price":    120,
    "status":   "active",
    "category": "electronics",
    "stock":    75,
    "products": []interface{}{
        map[string]interface{}{
            "info": map[string]interface{}{
                "title": "Smartphone",
                "available": true,
            },
        },
    },
    "discount": nil,
}
```

And here's how you might define a rule:

```go
group and
    price gt 100
    status eq "active"
    group or
        # Checking category and stock
        category in ["electronics", "furniture"]
        stock between 50 and 100
        # Product-specific checks
        products[0].info.title eq "Smartphone"
        products[0].info.available eq true
        discount eq nil
```

This rule checks:

- If price is greater than 100
- If status is "active"
- In the or group, it checks:
    - If category is either "electronics" or "furniture"
    - If stock is between 50 and 100
    - If the title of the first product is "Smartphone"
    - If the first product is available
    - If discount is nil



You can evaluate this rule with the following code:

```go
result, err := logix.EvaluateLogix(input, context)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Evaluation Result:", result)
}
```

You can also load the context from a JSON file like this:
```go
context, err := logix.LoadContextFromFile("context.json")
if err != nil {
    fmt.Println("Error loading context:", err)
    return
}

result, err := logix.EvaluateLogix(input, context)
if err != nil {
    fmt.Println("Error evaluating rule:", err)
} else {
    fmt.Println("Evaluation Result:", result)
}
```

## TODOs

- Improve error messages to show line and column numbers for better debugging
- Remove panics from the parser and return proper errors instead
