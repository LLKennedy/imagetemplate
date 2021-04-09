import React from 'react';
import { CanvasWrapper, RGBA } from './lib/canvas';

interface Props { }

class State {
  public ref: React.RefObject<HTMLCanvasElement> = React.createRef<HTMLCanvasElement>();
}

class App extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = new State();
  }
  render() {
    return <div>
      <header>
        Canvas Test
        </header>
      <canvas ref={this.state.ref} />
    </div>;
  }
  componentDidMount() {
    window.requestAnimationFrame(this.draw.bind(this));
  }
  async draw() {
    if (this.state.ref.current === undefined || this.state.ref.current === null) {
      window.requestAnimationFrame(this.draw.bind(this));
      return;
    }
    const ctx = this.state.ref.current.getContext("2d", {
      alpha: true,
      desynchronized: false,
    })
    if (ctx === null) {
      window.alert("canvas context was null");
      throw new Error("canvas context was null");
    }
    const wrapper = new CanvasWrapper(ctx);
    const red = new RGBA();
    red.R = 255;
    red.A = 127;
    await wrapper.Rectangle({ x: 0, y: 0 }, 100, 100, red);
  }
}

export default App;
