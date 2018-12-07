import React, { Component } from 'react'

class ButtonTable extends Component {
  state = {}
  render() {
    return (
      <div>
        <p>
          <button className="button" onClick={this.props.add}>
            Add
          </button>
        </p>
        <p>
          <button className="button" onClick={this.props.delete}>
            Remove selected
          </button>
        </p>
        <p>
          <button className="button" onClick={this.props.clear}>
            Clear
          </button>
        </p>
        <p>
          <a id="downloadAnchorElem" style={{ display: 'none' }} />
          <button className="button" onClick={this.props.export}>
            Exported To JSON
          </button>
        </p>
        <p>
          <button className="button" onClick={this.props.sortKey}>
            Sort By Name
          </button>
        </p>
        <p>
          <button className="button" onClick={this.props.sortValue}>
            Sort By Value
          </button>
        </p>
      </div>
    )
  }
}

export default ButtonTable
