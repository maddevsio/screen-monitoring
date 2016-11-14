import React, {PureComponent} from 'react';
import WidgetRow from './WidgetRow.jsx';

const WIDGETS_UNREG_URL = '/dashboard/v1/widgets/unregistered';
const TABLE_VIEW = 0;
const PAGES_PREVIEW = 1;

class UnregisteredWidgetList extends PureComponent {

  static defaultProps = {
    widgets: [],
    viewState: TABLE_VIEW,
    selected: null
  }

  static propTypes = {
    widgets: React.PropTypes.array
  }

  constructor(props) {
      super(props);
      this.state = {
        widgets: props.widgets,
        viewState: props.viewState,
        selected: null
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

  _onRegisterWidget = (w) => {

  }

  _renderTable = () => {
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
                  <th>&nbsp;</th>
                </tr>
              </thead>
              <tbody>
                  {
                    this.state.widgets.map((w, idx) => {
                      return (
                        <WidgetRow key={idx} widget={w} onClick={this._onRegisterWidget}/>
                      );
                    })
                  }
              </tbody>
            </table>
          </div>
      </div>
    )
  }

  _renderPages = () => {

  }

  render(){
    if (this.state.viewState == TABLE_VIEW) {
      return this._renderTable();
    } else if(this.state.viewState == PAGES_PREVIEW) {
      return this._renderPages();
    }
  }
}

export default UnregisteredWidgetList;
