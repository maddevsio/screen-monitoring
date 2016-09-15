import React, { PropTypes } from 'react';

class Preview extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
        id: props.id,
        content: props.content,
        width: props.width,
        height: props.height,
        url:    props.url
    }
  }
  componentWillReceiveProps(nextProps) {
        this.setState({
          content: nextProps.url
        });
        this.refreshIframe();
  }
  componentDidMount() {
      this.refreshIframe();
  }
  refreshIframe() {
    this.refs.iframe.src = this.state.url;
  }
  render() {
    return (
      <div>
        <iframe id={this.state.id}
                width={this.state.width}
                height={this.state.height}
                frameBorder="0"
                style={{float: 'left'}}
                ref="iframe"/>
      </div>
    );
  }
}

export default Preview;