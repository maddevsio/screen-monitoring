import React, {PureComponent} from 'react';

const WIDGETS_UNREG_URL = '/dashboard/v1/widgets/unregistered';

class UnregisteredWidgetList extends PureComponent {

  static defaultProps = {
    widgets: []
  }

  static propTypes = {
    widgets: React.PropTypes.array
  }

  constructor(props) {
      super(props);
      this.state = {
        widgets: props.widgets
      }
  }

  componentDidMount(){
    this._getUnregisteredWidgetsList();
  }

  _getUnregisteredWidgetsList(){
     var self = this;
     fetch(WIDGETS_UNREG_URL)
       .then((response) => response.json())
       .then(self._receiveResponse)
       .then(function (data){
          self.setState({widgets:data});
       })
       .catch((ex) => {
           console.error(ex);
       });
  }

  render(){
    return (
      <div className="panel panel-default" style={{"marginTop": "15px"}}>
          <div className="panel-heading">
            <b>Unregistered widgets:</b>
          </div>
          <div className="panel-body">
            <table className="table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Width</th>
                  <th>Height</th>
                  <th>Url</th>
                </tr>
              </thead>
              <tbody>
                  {
                    this.state.widgets.map((w, idx) => {
                      return (
                        <tr key={idx}>
                          <th scope="row">{w.id}</th>
                          <td>{w.width}</td>
                          <td>{w.height}</td>
                          <td>{w.url}</td>
                        </tr>
                      )
                    })
                  }
              </tbody>
            </table>
          </div>
      </div>
    )
  }
}

export default UnregisteredWidgetList;
