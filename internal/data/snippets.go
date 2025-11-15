package data

import "math/rand"

// CodeSnippet represents a code snippet for typing
type CodeSnippet struct {
	Language string
	Code     string
	Title    string
}

// GoSnippets contains Go code snippets
var GoSnippets = []CodeSnippet{
	{
		Language: "go",
		Title:    "HTTP Server",
		Code: `func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}`,
	},
	{
		Language: "go",
		Title:    "Error Handling",
		Code: `if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}`,
	},
	{
		Language: "go",
		Title:    "Goroutine",
		Code: `go func() {
    for msg := range ch {
        fmt.Println(msg)
    }
}()`,
	},
	{
		Language: "go",
		Title:    "Interface Implementation",
		Code: `type Reader interface {
    Read(p []byte) (n int, err error)
}`,
	},
	{
		Language: "go",
		Title:    "Struct Definition",
		Code: `type User struct {
    ID        int64
    Name      string
    Email     string
    CreatedAt time.Time
}`,
	},
	{
		Language: "go",
		Title:    "Method",
		Code: `func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("name required")
    }
    return nil
}`,
	},
}

// JSSnippets contains JavaScript code snippets
var JSSnippets = []CodeSnippet{
	{
		Language: "javascript",
		Title:    "Async Function",
		Code: `async function fetchData() {
    const response = await fetch(url);
    return response.json();
}`,
	},
	{
		Language: "javascript",
		Title:    "Arrow Function",
		Code: `const sum = (a, b) => a + b;
const square = x => x * x;`,
	},
	{
		Language: "javascript",
		Title:    "Promise",
		Code: `promise
    .then(result => console.log(result))
    .catch(error => console.error(error));`,
	},
	{
		Language: "javascript",
		Title:    "Destructuring",
		Code: `const { name, age } = user;
const [first, ...rest] = array;`,
	},
	{
		Language: "javascript",
		Title:    "Class",
		Code: `class Animal {
    constructor(name) {
        this.name = name;
    }
    speak() {
        console.log(this.name + ' makes a sound.');
    }
}`,
	},
	{
		Language: "javascript",
		Title:    "Map Function",
		Code: `const doubled = numbers.map(n => n * 2);
const filtered = items.filter(x => x > 10);`,
	},
}

// PythonSnippets contains Python code snippets
var PythonSnippets = []CodeSnippet{
	{
		Language: "python",
		Title:    "List Comprehension",
		Code: `squares = [x**2 for x in range(10)]
evens = [x for x in numbers if x % 2 == 0]`,
	},
	{
		Language: "python",
		Title:    "Function",
		Code: `def greet(name: str) -> str:
    return f"Hello, {name}!"`,
	},
	{
		Language: "python",
		Title:    "Class",
		Code: `class Person:
    def __init__(self, name: str, age: int):
        self.name = name
        self.age = age`,
	},
	{
		Language: "python",
		Title:    "Context Manager",
		Code: `with open('file.txt', 'r') as f:
    content = f.read()
    print(content)`,
	},
	{
		Language: "python",
		Title:    "Lambda",
		Code: `multiply = lambda x, y: x * y
squares = map(lambda x: x**2, range(10))`,
	},
	{
		Language: "python",
		Title:    "Decorator",
		Code: `@app.route('/users')
def get_users():
    return jsonify(users)`,
	},
}

// RustSnippets contains Rust code snippets
var RustSnippets = []CodeSnippet{
	{
		Language: "rust",
		Title:    "Match Expression",
		Code: `match result {
    Ok(value) => println!("Success: {}", value),
    Err(e) => eprintln!("Error: {}", e),
}`,
	},
	{
		Language: "rust",
		Title:    "Struct",
		Code: `struct User {
    username: String,
    email: String,
    active: bool,
}`,
	},
	{
		Language: "rust",
		Title:    "Implementation",
		Code: `impl User {
    fn new(username: String, email: String) -> Self {
        Self { username, email, active: true }
    }
}`,
	},
	{
		Language: "rust",
		Title:    "Option",
		Code: `fn find_user(id: u32) -> Option<User> {
    users.iter().find(|u| u.id == id).cloned()
}`,
	},
	{
		Language: "rust",
		Title:    "Iterator",
		Code: `let sum: i32 = numbers
    .iter()
    .filter(|&&x| x > 0)
    .sum();`,
	},
	{
		Language: "rust",
		Title:    "Trait",
		Code: `trait Summary {
    fn summarize(&self) -> String;
}`,
	},
}

// TypeScriptSnippets contains TypeScript code snippets
var TypeScriptSnippets = []CodeSnippet{
	{
		Language: "typescript",
		Title:    "Interface",
		Code: `interface User {
    id: number;
    name: string;
    email?: string;
}`,
	},
	{
		Language: "typescript",
		Title:    "Generic Function",
		Code: `function identity<T>(arg: T): T {
    return arg;
}`,
	},
	{
		Language: "typescript",
		Title:    "Type Guard",
		Code: `function isString(value: unknown): value is string {
    return typeof value === 'string';
}`,
	},
	{
		Language: "typescript",
		Title:    "Async/Await",
		Code: `async function fetchUser(id: number): Promise<User> {
    const response = await fetch('/api/users/' + id);
    return response.json();
}`,
	},
	{
		Language: "typescript",
		Title:    "Union Types",
		Code: `type Status = 'pending' | 'active' | 'inactive';
type Result = Success | Error;`,
	},
}

// CppSnippets contains C++ code snippets
var CppSnippets = []CodeSnippet{
	{
		Language: "cpp",
		Title:    "Vector Operations",
		Code: `std::vector<int> nums = {1, 2, 3, 4, 5};
for (const auto& n : nums) {
    std::cout << n << std::endl;
}`,
	},
	{
		Language: "cpp",
		Title:    "Class Definition",
		Code: `class Rectangle {
private:
    int width, height;
public:
    Rectangle(int w, int h) : width(w), height(h) {}
    int area() const { return width * height; }
};`,
	},
	{
		Language: "cpp",
		Title:    "Template",
		Code: `template<typename T>
T max(T a, T b) {
    return (a > b) ? a : b;
}`,
	},
	{
		Language: "cpp",
		Title:    "Smart Pointer",
		Code: `auto ptr = std::make_unique<Widget>();
std::shared_ptr<Data> data = std::make_shared<Data>();`,
	},
	{
		Language: "cpp",
		Title:    "Lambda",
		Code: `auto sum = [](int a, int b) { return a + b; };
std::sort(vec.begin(), vec.end(), [](int a, int b) { return a > b; });`,
	},
}

// JavaSnippets contains Java code snippets
var JavaSnippets = []CodeSnippet{
	{
		Language: "java",
		Title:    "Class",
		Code: `public class Person {
    private String name;
    private int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
}`,
	},
	{
		Language: "java",
		Title:    "Stream API",
		Code: `List<Integer> numbers = Arrays.asList(1, 2, 3, 4, 5);
int sum = numbers.stream()
    .filter(n -> n % 2 == 0)
    .mapToInt(Integer::intValue)
    .sum();`,
	},
	{
		Language: "java",
		Title:    "Optional",
		Code: `Optional<User> user = userRepository.findById(id);
return user.orElseThrow(() -> new NotFoundException());`,
	},
	{
		Language: "java",
		Title:    "Try-Catch",
		Code: `try {
    result = riskyOperation();
} catch (IOException e) {
    logger.error("Failed", e);
} finally {
    cleanup();
}`,
	},
	{
		Language: "java",
		Title:    "Interface",
		Code: `public interface Repository<T, ID> {
    Optional<T> findById(ID id);
    List<T> findAll();
    T save(T entity);
}`,
	},
}

// CSharpSnippets contains C# code snippets
var CSharpSnippets = []CodeSnippet{
	{
		Language: "csharp",
		Title:    "LINQ Query",
		Code: `var query = from user in users
            where user.Age > 18
            select user.Name;`,
	},
	{
		Language: "csharp",
		Title:    "Property",
		Code: `public class Person {
    public string Name { get; set; }
    public int Age { get; private set; }
}`,
	},
	{
		Language: "csharp",
		Title:    "Async Method",
		Code: `public async Task<User> GetUserAsync(int id) {
    var response = await httpClient.GetAsync($"/api/users/{id}");
    return await response.Content.ReadAsAsync<User>();
}`,
	},
	{
		Language: "csharp",
		Title:    "Lambda Expression",
		Code: `var evenNumbers = numbers.Where(n => n % 2 == 0);
var doubled = numbers.Select(n => n * 2);`,
	},
	{
		Language: "csharp",
		Title:    "Null Coalescing",
		Code: `string name = user?.Name ?? "Unknown";
int count = list?.Count ?? 0;`,
	},
}

// RubySnippets contains Ruby code snippets
var RubySnippets = []CodeSnippet{
	{
		Language: "ruby",
		Title:    "Class",
		Code: `class Person
  attr_accessor :name, :age
  
  def initialize(name, age)
    @name = name
    @age = age
  end
end`,
	},
	{
		Language: "ruby",
		Title:    "Block",
		Code: `numbers.each do |n|
  puts n * 2
end`,
	},
	{
		Language: "ruby",
		Title:    "Map",
		Code: `squared = numbers.map { |n| n ** 2 }
evens = numbers.select { |n| n.even? }`,
	},
	{
		Language: "ruby",
		Title:    "String Interpolation",
		Code: `name = "Alice"
greeting = "Hello, #{name}!"`,
	},
	{
		Language: "ruby",
		Title:    "Hash",
		Code: `user = {
  name: "John",
  age: 30,
  email: "john@example.com"
}`,
	},
}

// PHPSnippets contains PHP code snippets
var PHPSnippets = []CodeSnippet{
	{
		Language: "php",
		Title:    "Function",
		Code: `function greet($name) {
    return "Hello, " . $name . "!";
}`,
	},
	{
		Language: "php",
		Title:    "Class",
		Code: `class User {
    private $name;
    private $email;
    
    public function __construct($name, $email) {
        $this->name = $name;
        $this->email = $email;
    }
}`,
	},
	{
		Language: "php",
		Title:    "Array Map",
		Code: `$squared = array_map(function($n) {
    return $n * $n;
}, $numbers);`,
	},
	{
		Language: "php",
		Title:    "Null Coalescing",
		Code: `$username = $_GET['user'] ?? 'guest';
$config = $options['debug'] ?? false;`,
	},
	{
		Language: "php",
		Title:    "Arrow Function",
		Code: `$doubled = array_map(fn($n) => $n * 2, $numbers);
$filtered = array_filter($items, fn($x) => $x > 10);`,
	},
}

// GetCodeSnippet returns a random code snippet for the specified language
func GetCodeSnippet(language string) CodeSnippet {
	var snippets []CodeSnippet

	switch language {
	case "go":
		snippets = GoSnippets
	case "javascript", "js":
		snippets = JSSnippets
	case "python", "py":
		snippets = PythonSnippets
	case "rust", "rs":
		snippets = RustSnippets
	case "typescript", "ts":
		snippets = TypeScriptSnippets
	case "cpp", "c++":
		snippets = CppSnippets
	case "java":
		snippets = JavaSnippets
	case "csharp", "cs", "c#":
		snippets = CSharpSnippets
	case "ruby", "rb":
		snippets = RubySnippets
	case "php":
		snippets = PHPSnippets
	default:
		// Return random from any language
		allSnippets := append(GoSnippets, JSSnippets...)
		allSnippets = append(allSnippets, PythonSnippets...)
		allSnippets = append(allSnippets, RustSnippets...)
		allSnippets = append(allSnippets, TypeScriptSnippets...)
		allSnippets = append(allSnippets, CppSnippets...)
		allSnippets = append(allSnippets, JavaSnippets...)
		allSnippets = append(allSnippets, CSharpSnippets...)
		allSnippets = append(allSnippets, RubySnippets...)
		allSnippets = append(allSnippets, PHPSnippets...)
		snippets = allSnippets
	}

	if len(snippets) == 0 {
		return GoSnippets[0]
	}

	return snippets[rand.Intn(len(snippets))]
}
