<template>
  <div class="body_bg">
    <el-container class="outside">
      <el-header class="title">WebSocket即时在线聊天室
      </el-header>
      <el-container>
        <el-aside width="300px">
          <span class="text">用户</span>
          <el-input class="input" placeholder="请输入姓名" v-model="username" clearable></el-input>
          <el-button type="danger" @click="exit">退出聊天室</el-button>
          <el-button type="primary" @click="join">加入聊天室</el-button>
          <el-input class="input"
                    placeholder="请输入内容"
                    v-model="input"
                    clearable
                    type="textarea"
                    :rows="20">
          </el-input>
          <el-button type="warning" class="send_button" @click="sendAll">发送消息</el-button>
        </el-aside>
        <el-main>
          <span class="text">群聊</span>
          <el-card class="message">
            <el-input
                type="textarea"
                :autosize="{ minRows: 26, maxRows: 26}"
                v-model="message"
                disabled
                class="message_text">
            </el-input>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data() {
    return {
      input: '',
      message: '',
      username: '',
      urlPrefix: 'ws://localhost:8080/chat-room/',
      ws: null,
    }
  },
  methods: {
    join() {
      this.url = this.urlPrefix + this.username
      this.ws = new WebSocket(this.url)
      this.wsOpen()
      this.wsOnmessage()
      this.ws.onclose = function () {
        this.message += '用户[' + this.username + '] 已经离开聊天室' + '\n'
        console.log("关闭websocket连接")
      }
    },
    wsOpen() {
      this.ws.onopen = function () {
        console.log('connected!')
      }
    },
    wsOnmessage() {
      let _this = this
      this.ws.onmessage = function (event) {
        _this.message += event.data + '\n'
      }
    },
    sendAll() {
      if(this.ws) {
        this.ws.send(this.input)
        this.input = ''
        this.wsOnmessage()
      }

    },
    exit() {
      if(this.ws) {
        this.ws.close()
      }
    },

  }
}
</script>

<style lang="less" scoped>
.body_bg {
  position: absolute;
  height: 100%;
  width: 100%;
  top:0;
  left: 0;
  overflow-y: auto;
  background-color: #fff;
}

.outside {
  height: 100%;
}

.title {
  background-color: #083aa1;
  color: white;
  height: 60px;
  width: 100%;
  font-size: 28px;
  text-align: center;
  line-height: 60px;
}

.el-main {
  background-color: aliceblue;
  height: 100%;
  width: 100%;
}

.el-aside {
  background-color: #2c3e50;
  color: white;
}

.input {
  padding: 10%;
  width: 80%;
  font-size: 20px;
  font-family: 'Microsoft YaHei';
}

.text {
  position: relative;
  top: 10px;
  font-size: 24px;
  font-family: 'Microsoft YaHei';
  font-weight: bold;
}

.message {
  position: relative;
  top: 40px;
  height: 93%;
}

.send_button {
  position: relative;
  top: -10px;
  font-size: 18px;
  font-family: 'Microsoft YaHei';
}

.message_text {
  position: relative;
  top: -5px;
  font-size: 22px;
  font-weight: bold;
}
</style>
