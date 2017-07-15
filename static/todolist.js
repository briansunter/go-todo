var TodoList = React.createClass({
  loadTasksFromServer : function(){
    $.ajax({
      url: this.props.listUrl,
      type: 'GET',
      dataType: 'json',
      success: function (data, textStatus, jqXHR) {
        this.setState({tasks:data});
        console.log(data);
      }.bind(this),
      error: function (jqXHR, textStatus, errorThrown) {
        console.error(errorThrown);
      }
    });
  },
  getInitialState : function(){
    return {tasks : []};
  },
  componentDidMount: function() {
    this.loadTasksFromServer();
  },
  render: function(){
    var todoNodes = this.state.tasks.map(function(task){
      return (
      <TodoTask task={task}></TodoTask>
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
  <TodoList listName="Test Todo List" listUrl="http://localhost:8080/tasks"/>,
  document.getElementById('content')
);
