import React, {PureComponent} from 'react';

const css = {
  page: {
    "width": "150px",
    "height": "150px",
    "float":"left",
    "margin": "5px"
  }
}

class PageButton extends PureComponent {

  static defaultProps = {
    page: {id:0, title:"No page"}
  }

  static propTypes = {
    page: React.PropTypes.object.isRequired,
    onClick: React.PropTypes.func
  }

  constructor(props) {
    super(props);
    this.state = {
      page: props.page
    }
  }

  _renderPageVisibility = (visibility) => {
    if (!visibility) {
      return (
          <span className="label label-warning">Invisible page</span>
      )
    }
    return (<span></span>);
  }

  _clickPage = (e) => {
    if (this.props.onClick) {
      this.props.onClick(this.state.page);
    }
    return false;
  }

  _getWidgetsCount = () => {
    var page = this.state.page;
    return (page && page.content ? page.content.length : 0);
  }

  render(){
    var p = this.state.page;
    return (
      <div className="btn btn-default"
        style={css.page}
        onClick={this._clickPage}>
        <div>{p.title}</div>
        <div>{this._renderPageVisibility(p.visible)}</div>
        <div>Widgets: <strong>{this._getWidgetsCount()}</strong></div>
      </div>
    );
  }
}


export default PageButton;
