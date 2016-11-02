import React, { PropTypes, PureComponent } from 'react';

class Preview extends PureComponent {
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
  refreshIframe = () => {
    this._frame.src = this.state.url;
  }
  render() {
    return (
      <div>
        <iframe id={this.state.id}
                width={this.state.width}
                height={this.state.height}
                frameBorder="0"
                src={this.state.url}
                style={{float: 'left'}}
                ref={(c) => this._frame = c}/>
      </div>
    );
  }
}

export default Preview;
