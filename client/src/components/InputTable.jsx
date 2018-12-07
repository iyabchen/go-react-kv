import React, { Component } from 'react'

class InputTable extends Component {
  constructor(props) {
    super(props)
    const numRows = this.props.numRows
    this.inputs = []
    for (let i = 0; i < numRows; i++) {
      this.inputs.push(React.createRef())
    }
  }
  render() {
    let rows = []
    const numRows = this.props.numRows
    for (let i = 0; i < numRows; i++) {
      let color = 'white'
      if (this.props.selectedInputID === i.toString()) {
        color = 'blue'
      }
      rows.push(
        <InputRow
          color={color}
          onClick={this.props.onClick}
          onChange={this.props.onChange}
          id={i}
          key={'input' + i}
        />
      )
    }
    return (
      <table width="300px" border="1">
        <tbody>{rows}</tbody>
      </table>
    )
  }
}

function InputRow(props) {
  const input = (
    <input
      id={props.id}
      onChange={props.onChange}
      style={{ width: '200px' }}
      type="text"
    />
  )

  return (
    <tr>
      <td
        id={props.id}
        style={{ backgroundColor: props.color }}
        key={props.id}
        onClick={props.onClick}>
        {input}
      </td>
    </tr>
  )
}
export default InputTable
