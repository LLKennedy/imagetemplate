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
      <canvas ref={this.state.ref} width={1000} height={1000} />
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
    await wrapper.Rectangle({ x: 0, y: 0 }, 100, 50, new RGBA(255, 0, 0, 127));
    await wrapper.Rectangle({ x: 0, y: 0 }, 50, 100, new RGBA(255, 255, 0, 127));
    await wrapper.Circle({ x: 25, y: 25 }, 10, new RGBA(0, 0, 255, 127));
  }
}

export default App;
