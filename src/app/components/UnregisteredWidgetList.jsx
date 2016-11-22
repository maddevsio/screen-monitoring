import React, {PureComponent} from 'react';
import WidgetRow from './WidgetRow.jsx';
import PageButton from './PageButton.jsx';
import { hashHistory } from 'react-router';
import AddPageForm from './AddPageForm.jsx';

const WIDGETS_UNREG_URL = '/dashboard/v1/widgets/unregistered';
const GET_PAGES_URL = '/dashboard/v1/pages';
const TABLE_VIEW = 0;
const PAGES_PREVIEW = 1;
const PAGES_ADD = 2;

const css = {
  page: {
    "width": "150px",
    "height": "150px",
    "float":"left",
    "display": "flex",
    "alignItems": "center",
    "justifyContent": "center",
    "margin": "5px"
  }
}

class UnregisteredWidgetList extends PureComponent {

  static defaultProps = {
    widgets: [],
    pages: [],
    viewState: TABLE_VIEW,
    selected: null
  }

  static propTypes = {
    widgets: React.PropTypes.array
  }

  constructor(props) {
      super(props);
      this.state = {
        pages: props.pages,
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
       .then(r => r.json())
       .then(function (data){
          self.setState({widgets:data});
       })
       .catch((ex) => {
           console.error(ex);
       });
  }

  _onRegisterWidget = (w) => {
    var selectedId = w.id;
    hashHistory.push({
      pathname: '/pages/list',
      state: { selectedWidgetId: selectedId }
    });
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
    );
  }

  render() {
      return this._renderTable();
  }
}

export default UnregisteredWidgetList;
