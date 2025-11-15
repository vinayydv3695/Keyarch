class Person:
    def __init__(self, name: str, age: int):
        self.name = name
        self.age = age
    
    def greet(self) -> str:
        return f"Hello, I'm {self.name} and I'm {self.age} years old"
    
    def celebrate_birthday(self):
        self.age += 1
        print(f"Happy birthday! Now {self.age} years old")

person = Person("Alice", 25)
print(person.greet())
