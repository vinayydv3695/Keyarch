class Animal {
    constructor(name, species) {
        this.name = name;
        this.species = species;
    }

    makeSound() {
        console.log(`${this.name} makes a sound`);
    }

    introduce() {
        return `I am ${this.name}, a ${this.species}`;
    }
}

const dog = new Animal('Buddy', 'Dog');
dog.makeSound();
