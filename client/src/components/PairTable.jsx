import React, { Component } from 'react'

class PairTable extends Component {
  state = {}

  render() {
    let rows = []
    let numRows = this.props.numRows
    let onClick = this.props.onClick
    let seletedPairID = this.props.selectedPairID

    let pairs = this.props.pairs
    let i = 0
    for (let p of pairs) {
      let color = 'white'
      if (p.id === seletedPairID) {
        color = 'blue'
      }
      rows.push(
        <Pair
          onClick={onClick}
          key={'kv' + i}
          id={p.id}
          pkey={p.key}
          pvalue={p.value}
          color={color}
        />
      )
      i++
    }
    for (; i < numRows; i++) {
      rows.push(
        <Pair onClick={onClick} key={'kv' + i} id="" pkey="" pvalue="" />
      )
    }

    return (
      <table width="300px" border="1">
        <tbody>{rows}</tbody>
      </table>
    )
  }
}

function Pair(props) {
  const text =
    props.pkey === '' ? (
      ' '
    ) : (
      <span>
        {props.pkey} = {props.pvalue}
      </span>
    )

  return (
    <tr>
      <td
        style={{ backgroundColor: props.color }}
        id={props.id}
        onClick={props.onClick}>
        {text}
      </td>
    </tr>
  )
}
export default PairTable
