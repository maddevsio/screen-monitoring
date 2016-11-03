import React, {PureComponent} from 'react';

class WidgetRow extends PureComponent {

  static propTypes = {
    widget: React.PropTypes.object,
    onClick: React.PropTypes.func
  }

  constructor(props) {
    super(props);
    this.state = {
      widget: props.widget
    }
  }

  _onClick = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (this.props.onClick) {
      this.props.onClick(this.state.widget);
    }
  }

  render(){
    var w = this.state.widget;
    return (
      <tr>
        <th scope="row">{w.id}</th>
        <td>{w.width}</td>
        <td>{w.height}</td>
        <td>{w.url}</td>
        <td>
          <a onClick={this._onClick} className="btn btn-sm btn-success"
             href="javascript:void(0)">Register to page</a>
        </td>
      </tr>
    );
  }
}

export default WidgetRow;
