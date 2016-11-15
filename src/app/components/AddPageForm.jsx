import React, {PureComponent} from 'react';
import serialize from 'form-serialize';

class AddPageForm extends PureComponent {

  static propTypes = {
    page: React.PropTypes.object,
    onCreate: React.PropTypes.func,
    onCancel: React.PropTypes.func
  }

  static defaultProps = {
    page: {title: "", visible: true}
  }

  constructor(props) {
    super(props);
    this.state = {
        page: props.page
    }
  }

  _pageTitleChanged = (e) => {
      var page = this.state.page;
      page.title = e.target.value;
      this.setState({page: page});
  }

  _pageVisibleChanged = (e) => {
      var page = this.state.page;
      page.visible = e.target.checked;
      this.setState({page: page});
  }


  _normalizePageData = (p) => {
    if (p.visible) {
      p.visible = (p.visible == "true");
    }
    return p;
  }

  _onCreate = (e) => {
    e.preventDefault();
    var data = this._normalizePageData(
      serialize(e.target, { hash: true })
    );
    if (this.props.onCreate) {
      this.props.onCreate(data);
    }
  }

  _onCancel = (e) => {
    e.preventDefault();
    if (this.props.onCancel) {
      this.props.onCancel();
    }
  }

  render() {
    var page = this.state.page;
    return (
      <div className="well well-sm col-md-6 col-md-offset-3" style={{"marginTop": "30px"}}>
        <h3>New page</h3>
        <form className="centered" onSubmit={this._onCreate}>
            <div className="form-group row">
              <div className="col-sm-12">
                  <label htmlFor="page_title">Page title</label>
                  <input
                    id="page_title"
                    name="title"
                    type="text"
                    className="form-control"
                    placeholder="Page title"
                    onChange={this._pageTitleChanged}
                    defaultValue={page.title}/>
              </div>
            </div>
            <div className="form-group row">
              <div className="col-sm-12">
                <div className="form-check">
                  <label htmlFor="page_visible" className="form-check-label">
                    Visible? &nbsp;&nbsp;
                  </label>
                  <input id="page_visible"
                      name="visible"
                      className="form-check-input" type="checkbox"
                      onChange={this._pageVisibleChanged}
                      value={page.visible}
                      defaultChecked={page.visible}/>
                </div>
              </div>
            </div>
            <div className="form-group row">
              <div className="offset-sm-2 col-sm-12">
                <button type="submit" className="btn btn-primary">
                  Create
                </button>
                <button onClick={this._onCancel} className="btn btn-danger pull-right">
                  Cancel
                </button>
              </div>
            </div>
          </form>
        </div>
    );
  }
}

export default AddPageForm;
