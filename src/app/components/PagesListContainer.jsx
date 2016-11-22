import React, {PureComponent} from 'react';
import { hashHistory } from 'react-router';
import PageButton from './PageButton.jsx';

const GET_PAGES_URL = '/dashboard/v1/pages';

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

class PagesListContainer extends PureComponent {

  static defaultProps = {
    pages: []
  }

  constructor(props) {
    super(props);
    this.state = {
      pages: [],
      loading: false
    }
  }

  componentDidMount() {
    this._queryAllPages();
  }

  _queryAllPages = () => {
    this.setState({loading: true});
    return fetch(GET_PAGES_URL)
      .then(r => r.json())
      .then(data => {
         this.setState({
           loading: false,
           pages: data
         });
      })
      .catch(ex => {
          console.error(ex);
          this.setState({loading: false});
      });
  }

  _showAddPageForm = () => {
    let selectedWidgetId = this.props.location.state.selectedWidgetId;
    hashHistory.push({
       pathname: '/pages/new',
       state: { selectedWidgetId: selectedWidgetId }
    });
  }

  _postRegisterWidgetOnPage(pageId, widgetId){
    const REGISTER_WIDGET_ON_PAGE_URL = `/dashboard/v1/register/${widgetId}/page/${pageId}`;
    fetch(REGISTER_WIDGET_ON_PAGE_URL)
      .then(r => r.json())
      .then(data => {
         if (data.Success) {
          this._queryAllPages().then(() => {
                hashHistory.push({
                   pathname: '/pages/list',
                   state: { selectedWidgetId: null }
                });
          });
         } else {
           alert(data.error);
         }
      })
      .catch((ex) => {
          console.error(ex);
      });
  }

  _choosePage = (p) => {
    var selected = this.props.location.state.selectedWidgetId;
    if (selected) {
      this._postRegisterWidgetOnPage(p.id, selected);
    }
  }

  _onGoToUnregistredList = (e) => {
    e.preventDefault();
    e.stopPropagation();
    hashHistory.push("/unregistered");
  }

  render () {
    if (this.state.loading) {
      return (
        <div className="panel panel-default" style={{"marginTop": "15px"}}>
            <div className="panel-heading">
              <b>Pages list:</b>
            </div>
            <div className="panel-body">
              <div className="alert alert-info">
                  <strong>Please wait!</strong><span>  Loading pages list...</span>
              </div>
            </div>
        </div>
      )
    } else {
      return (
        <div className="panel panel-default" style={{"marginTop": "15px"}}>
            <div className="panel-heading">
              <div className="row">
                 <div className="col-xs-10">
                   <b>Pages list:</b>
                 </div>
                 <div className="col-xs-2">
                   <button onClick={this._onGoToUnregistredList}
                           className="btn btn-sm btn-default pull-right">Unregistered widgets</button>
                 </div>
              </div>
            </div>
            <div className="panel-body">
                <div key={0} className="btn btn-success" onClick={this._showAddPageForm} style={css.page}>
                  <div><i className="fa fa-2x fa-plus" aria-hidden="true"></i></div>
                </div>
                {
                   this.state.pages.map(p => (<PageButton onClick={this._choosePage} key={p.id} page={p}/>))
                }
            </div>
        </div>
      )
    }
  }
}

export default PagesListContainer;
