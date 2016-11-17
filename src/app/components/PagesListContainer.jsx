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
    fetch(GET_PAGES_URL)
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
    hashHistory.push('/pages/new');
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
              <b>Pages list:</b>
            </div>
            <div className="panel-body">
                <div key={0} className="btn btn-success" onClick={this._showAddPageForm} style={css.page}>
                  <div><i className="fa fa-2x fa-plus" aria-hidden="true"></i></div>
                </div>
                {
                   this.state.pages.map(p => (<PageButton key={p.id} page={p}/>))
                }
            </div>
        </div>
      )
    }
  }
}

export default PagesListContainer;
