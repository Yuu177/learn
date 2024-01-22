# caffe

caffe 是一个清晰而高效的深度学习框架，是纯粹的 C++/CUDA 架构。

caffe 的全称是 Convolutional Architecture for Fast Feature Embedding（译为：快速特征嵌入的卷积体系结构），核心语言是 C++。caffe 的基本工作流程是设计建立在神经网络的一个简单假设，所有的计算都是层的形式表示的，网络层所做的事情就是输入数据，然后输出计算结果。比如卷积就是输入一幅图像，然后和这一层的参数（filter）做卷积，最终输出卷积结果。每层需要两种函数计算，一种是forward，从输入计算到输出；另一种是 backward，从上层给的 gradient 来计算相对于输入层的 gradient。这两个函数实现之后，我们就可以把许多层连接成一个网络，这个网络输入数据（图像，语音或其他原始数据），然后计算需要的输出（比如识别的标签）。在训练的时候，可以根据已有的标签计算 loss 和 gradient，然后用 gradient 来更新网络中的参数。

// TODO
