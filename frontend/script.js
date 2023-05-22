todoInputElement = null;
todoAddButtonElement = null;
todoListParentElement = null;
todos = []
document.addEventListener("DOMContentLoaded", function() {
    todoInputElement = document.getElementById("todo-input");
    todoAddButtonElement = document.getElementById("todo-add-button");
    todoListParentElement = document.getElementById("todo-list-parent");

    loadTodos();
});

async function addTodo() {
    todoName = todoInputElement.value;
    if (todoName == "") {
        return;
    }

    reqJSON = {"title": todoName, "priority": 0};

    ret = await fetch("http://" + window.location.hostname + ":3000/todos", {
        method: "POST",
        body: JSON.stringify(reqJSON),
    })

    if (ret.status != 201) {
        alert("Failed to add todo");
        return;
    }

    todoInputElement.value = "";
    todoAddButtonElement.disabled = true;
    await loadTodos();
}

async function loadTodos() {
    todos = await fetch("http://" + window.location.hostname + ":3000/todos")
    todos = await todos.json();

    
    todoListParentElement.innerHTML = "";
    for (todo of todos) {
        to_add = `
        <div class="todo-item">
            <div class="todo-item-title">${todo.title}</div>
            <div class="todo-item-delete" onclick="deleteTodo(${todo.id})">X</div>
        </div>
        `;

        if (todo.completed) {
            to_add = to_add.replace("todo-item", "todo-item todo-item-completed");
        }


        todoListParentElement.innerHTML += to_add;
    }

    todoItems = document.getElementsByClassName("todo-item");
    for (let i = 0; i < todoItems.length; i++) {
        todoItems[i].addEventListener("click", async function(event) {
            if (event.target.classList.contains("todo-item-delete")) {
                return;
            }

            await toggleTodoCompleted(todos[i].id, i);
        });
    }
}

async function toggleTodoCompleted(id, index) {
    res = await fetch("http://" + window.location.hostname + ":3000/todos/" + id, {
        method: "PATCH",
        body: JSON.stringify({"completed": !todos[index].completed}),    
    })

    if (res.status != 200) {
        alert("Failed to toggle todo completed");
        return;
    }

    await loadTodos();
}

function onTodoNameKeyUp(event) {
    // check length
    if (todoInputElement.value.length > 0) {
        todoAddButtonElement.disabled = false;
    } else {
        todoAddButtonElement.disabled = true;
        return;
    }
    
    if (event.keyCode == 13) {
        addTodo();
    }
}

async function deleteTodo(id) {
    res = await fetch("http://" + window.location.hostname + ":3000/todos/" + id, {
        method: "DELETE",
    })

    if (res.status != 200) {
        alert("Failed to delete todo");
        return;
    }

    await loadTodos();
}