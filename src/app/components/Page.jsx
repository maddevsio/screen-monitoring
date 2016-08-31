import React from 'react';

class Page extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            id: props.id,
            width: props.width,
            height: props.height,
            content: props.content
        };
    }

    componentWillReceiveProps(nextProps) {
          this.setState({
            content: nextProps.content
          });
    }

    render() {
        return (
            <div style={{
                          width: this.state.width + 'px',
                          height: this.state.height + 'px',
                          float: 'left' }}
                 dangerouslySetInnerHTML={{__html: this.state.content}}>
            </div>
        );
    }

}

export default Page;