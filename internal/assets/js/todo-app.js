class TodoApp {
    constructor(config) {
        this.config = config;
        this.todos = [];
        this.initializeEventListeners();
        this.loadTodos();
    }

    initializeEventListeners() {
        // Add todo form
        document.getElementById('addTodoForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.addTodo(new FormData(e.target));
        });

        // Search and filters
        document.getElementById('search').addEventListener('input', (e) => {
            this.filterTodos();
        });

        document.getElementById('category').addEventListener('change', () => {
            this.filterTodos();
        });

        document.getElementById('priority').addEventListener('change', () => {
            this.filterTodos();
        });

        // Archive toggle
        document.getElementById('showArchived').addEventListener('click', () => {
            this.toggleArchived();
        });
    }

    async loadTodos() {
        try {
            const response = await fetch(this.config.apiEndpoint);
            this.todos = await response.json();
            this.renderTodos();
        } catch (error) {
            console.error('Error loading todos:', error);
        }
    }

    async addTodo(formData) {
        try {
            const response = await fetch(this.config.apiEndpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(Object.fromEntries(formData)),
            });

            if (response.ok) {
                const todo = await response.json();
                this.todos.push(todo);
                this.renderTodos();
                this.config.onTodoUpdate?.(todo);
            }
        } catch (error) {
            console.error('Error adding todo:', error);
        }
    }

    async updateTodo(id, updates) {
        try {
            const response = await fetch(`${this.config.apiEndpoint}/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(updates),
            });

            if (response.ok) {
                const updatedTodo = await response.json();
                this.todos = this.todos.map(todo => 
                    todo.id === id ? updatedTodo : todo
                );
                this.renderTodos();
                this.config.onTodoUpdate?.(updatedTodo);
            }
        } catch (error) {
            console.error('Error updating todo:', error);
        }
    }

    async addComment(todoId, content) {
        try {
            const response = await fetch(`${this.config.apiEndpoint}/${todoId}/comments`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content }),
            });

            if (response.ok) {
                const updatedTodo = await response.json();
                this.todos = this.todos.map(todo => 
                    todo.id === todoId ? updatedTodo : todo
                );
                this.renderTodos();
            }
        } catch (error) {
            console.error('Error adding comment:', error);
        }
    }

    filterTodos() {
        const searchTerm = document.getElementById('search').value.toLowerCase();
        const category = document.getElementById('category').value;
        const priority = document.getElementById('priority').value;

        const filteredTodos = this.todos.filter(todo => {
            const matchesSearch = todo.description.toLowerCase().includes(searchTerm);
            const matchesCategory = !category || todo.category === category;
            const matchesPriority = !priority || todo.priority === parseInt(priority);
            return matchesSearch && matchesCategory && matchesPriority;
        });

        this.renderTodos(filteredTodos);
    }

    renderTodos(todosToRender = this.todos) {
        const todoList = document.getElementById('todoList');
        todoList.innerHTML = todosToRender.map(todo => this.renderTodoItem(todo)).join('');
    }

    renderTodoItem(todo) {
        return `
            <div class="todo-item bg-white p-4 rounded-lg shadow" data-id="${todo.id}">
                <!-- Todo item template -->
            </div>
        `;
    }
} 