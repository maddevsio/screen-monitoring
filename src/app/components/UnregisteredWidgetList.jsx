import React, {PureComponent} from 'react';
import WidgetRow from './WidgetRow.jsx';
import PageButton from './PageButton.jsx';
import { hashHistory } from 'react-router'

const WIDGETS_UNREG_URL = '/dashboard/v1/widgets/unregistered';
const GET_PAGES_URL = '/dashboard/v1/pages';
const TABLE_VIEW = 0;
const PAGES_PREVIEW = 1;

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
        widgets: props.widgets,
        viewState: props.viewState,
        selected: null
      }
  }

  componentDidMount(){
    this._getUnregisteredWidgetsList();
    this._getPagesList();
  }

  _getPagesList() {
    var self = this;
    fetch(GET_PAGES_URL)
      .then(r => r.json())
      .then(function (data){
         self.setState({pages:data});
      })
      .catch((ex) => {
          console.error(ex);
      });
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

  _postRegisterWidgetOnPage(pageId, widgetId){
    var self = this;
    const REGISTER_WIDGET_ON_PAGE_URL = `/dashboard/v1/register/${widgetId}/page/${pageId}`;
    fetch(REGISTER_WIDGET_ON_PAGE_URL)
      .then(r => r.json())
      .then(data => {
         if (data.Success) {
           hashHistory.push('/');
         } else {
           alert(data.error);
         }
      })
      .catch((ex) => {
          console.error(ex);
      });
  }

  _onRegisterWidget = (w) => {
      this.setState({
        selected: w,
        viewState: PAGES_PREVIEW
      });
  }

  _choosePage = (p) => {
    console.log("PAGE: ", p);
    var selected = this.state.selected;
    if (selected) {
      this._postRegisterWidgetOnPage(p.id, selected.id);
    }
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

  _renderPages = () => {
    return (
        <div className="panel panel-default" style={{"marginTop": "15px"}}>
            <div className="panel-heading">
              <b>Pages list:</b>
            </div>
            <div className="panel-body">
                <div key={0} className="btn btn-success" style={css.page}>
                  <div><i className="fa fa-2x fa-plus" aria-hidden="true"></i></div>
                </div>
                {
                   this.state.pages.map(p => (<PageButton onClick={this._choosePage} key={p.id} page={p}/>))
                }
            </div>
        </div>
    );
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
