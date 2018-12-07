import React, { Component } from 'react'
import InputTable from './InputTable'
import PairTable from './PairTable'
import ButtonTable from './ButtonTable'
import axios from 'axios'

const ADDR = process.env.REACT_APP_API_ADDR || ''

class Wrapper extends Component {
  constructor(props) {
    super(props)
    axios.get(ADDR + '/pair').then(response => {
      var pairs = response.data.slice(0)
      this.setState({ pairs: pairs })
    })
  }
  state = {
    selectedPairID: '',
    selectedInputID: '',
    selectedInput: '',
    pairs: []
  }
  render() {
    return (
      <div>
        <div className="left">
          <InputTable
            className="left"
            numRows={10}
            onClick={this.handleInputTableClick}
            onChange={this.handleInputChange}
            selectedInputID={this.state.selectedInputID}
          />
        </div>
        <div className="center">
          <ButtonTable
            className="center"
            delete={this.handleDelete}
            add={this.handleAdd}
            sortKey={this.handleSortByKey}
            sortValue={this.handleSortByValue}
            clear={this.handleClear}
            export={this.handleDownload}
          />
        </div>
        <div className="right">
          <PairTable
            className="right"
            numRows={10}
            pairs={this.state.pairs}
            onClick={this.handlePairTableClick}
            selectedPairID={this.state.selectedPairID}
          />
        </div>
      </div>
    )
  }

  handleInputTableClick = ev => {
    if (ev.target === ev.currentTarget) {
      this.setState({
        selectedInputID: ev.target.id,
        selectedInput: ev.target.children[0].value
      })
    }
  }

  handleInputChange = ev => {
    if (ev.target.id === this.state.selectedInputID) {
      this.setState({
        selectedInput: ev.target.value
      })
    }
  }
  handlePairTableClick = ev => {
    if (ev.target === ev.currentTarget) {
      this.setState({ selectedPairID: ev.target.id })
    }
  }

  handleAdd = () => {
    const str = this.state.selectedInput
    if (str === '') {
      alert('Please choose an input')
      return
    }
    let arr = str.split('=')
    if (arr.length !== 2) {
      alert('Must be separated by =')
      return
    }
    let pkey = arr[0].trim()
    let pvalue = arr[1].trim()
    if (!isAlphanumeric(pkey)) {
      alert('Key is not alphanumeric')
      return
    }
    if (!isAlphanumeric(pvalue)) {
      alert('Value is not alphanumeric')
      return
    }

    const pairs = this.state.pairs

    axios
      .post(ADDR + '/pair', {
        key: pkey,
        value: pvalue
      })
      .then(response => {
        console.log(response)
        let pair = response.data
        let newPairs = [...pairs, pair]
        console.log(newPairs)
        this.setState({ pairs: newPairs })
      })
  }

  handleDelete = () => {
    const selectedPairID = this.state.selectedPairID
    if (selectedPairID === '') {
      alert('Please choose a pair to remove')
      return
    }
    const pairs = this.state.pairs
    let i = 0
    for (; i < pairs.length; i++) {
      if (selectedPairID === pairs[i].id) {
        break
      }
    }
    let newPairs = [...pairs.slice(0, i), ...pairs.slice(i + 1)]
    axios.delete(ADDR + '/pair/' + selectedPairID).then(response => {
      this.setState({ pairs: newPairs })
    })
  }

  handleClear = () => {
    axios.get(ADDR + '/reset').then(response => {
      this.setState({ pairs: [] })
    })
  }

  handleSortByKey = () => {
    const pairs = this.state.pairs
    let newPairs = pairs.sort((a, b) => {
      let ca = a.key,
        cb = b.key
      if (ca < cb) return -1
      if (ca > cb) return 1
      return 0
    })
    this.setState({ pairs: newPairs })
  }

  handleSortByValue = () => {
    const pairs = this.state.pairs
    let newPairs = pairs.sort((a, b) => {
      let ca = a.value,
        cb = b.value
      if (ca < cb) return -1
      if (ca > cb) return 1
      return 0
    })
    this.setState({ pairs: newPairs })
  }

  handleDownload = () => {
    const pairs = this.state.pairs
    var dataStr =
      'data:text/json;charset=utf-8,' +
      encodeURIComponent(JSON.stringify(pairs))
    var dlAnchorElem = document.getElementById('downloadAnchorElem')
    dlAnchorElem.setAttribute('href', dataStr)
    dlAnchorElem.setAttribute('download', 'kv.json')
    dlAnchorElem.click()
  }
}

function isAlphanumeric(str) {
  const regexp = /[a-zA-Z0-9]+/
  return str.match(regexp)
}

export default Wrapper
