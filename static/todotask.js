var TodoTodo = React.createClass({
  onChange : function(){
    var self = this;
$.ajax({
  url: 'http://localhost:8080/todos/'+self.props.todo.id+ '/complete',
  type: 'PUT',
  dataType: 'json',
  success: function (data, textStatus, jqXHR) {
    console.log(data);

  },
  error: function (jqXHR, textStatus, errorThrown) {
    console.error(errorThrown);
  }
});

  },
  shouldCheck :function(){
    var completedDate = new Date(this.props.todo.completion_date);
    var referenceDate = new Date('1980-05-23')
    var hasBeenCompleted = (completedDate>referenceDate );
    console.log(hasBeenCompleted);
    return hasBeenCompleted;

  },
  render: function(){
    return(
      <li className="todoTodo">
      <h3>{this.props.todo.name}</h3>
      <h4>{this.props.todo.description}</h4>
      <h4>{this.props.todo.due_date}</h4>
      <input name="completed" type="checkbox" checked={this.shouldCheck()} onChange={this.onChange} />
      </li>
    );
  }

});
