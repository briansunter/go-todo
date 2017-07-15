var TodoList = React.createClass({
  loadTodosFromServer : function(){
    $.ajax({
      url: this.props.listUrl,
      type: 'GET',
      dataType: 'json',
      success: function (data, textStatus, jqXHR) {
        this.setState({todos:data});
        console.log(data);
      }.bind(this),
      error: function (jqXHR, textStatus, errorThrown) {
        console.error(errorThrown);
      }
    });
  },
  getInitialState : function(){
    return {todos : []};
  },
  componentDidMount: function() {
    this.loadTodosFromServer();
  },
  render: function(){
    var todoNodes = this.state.todos.map(function(todo){
      return (
      <TodoTodo todo={todo}></TodoTodo>
      );
    });

    return (
      <ul className="todoList">
      {todoNodes}
      </ul>
    );
  }
});

React.render(
  <TodoList listName="Test Todo List" listUrl="http://localhost:8080/todos"/>,
  document.getElementById('content')
);
