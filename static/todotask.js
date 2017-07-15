var TodoTask = React.createClass({
  onChange : function(){
    var self = this;
$.ajax({
  url: 'http://localhost:8080/tasks/'+self.props.task.id+ '/complete',
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
    var completedDate = new Date(this.props.task.completion_date);
    var referenceDate = new Date('1980-05-23')
    var hasBeenCompleted = (completedDate>referenceDate );
    console.log(hasBeenCompleted);
    return hasBeenCompleted;
  
  },
  render: function(){
    return(
      <li className="todoTask">
      <h3>{this.props.task.name}</h3>
      <h4>{this.props.task.description}</h4>
      <h4>{this.props.task.due_date}</h4>
      <input name="completed" type="checkbox" checked={this.shouldCheck()} onChange={this.onChange} />
      </li>
    );
  }

});
