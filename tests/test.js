function test() {
    console.log('hello world\n');
    return 'Hello Golang~';
}

function bar() {
    return self.foo();
}

module.exports = {
    test: test,
    foo: test,
};
