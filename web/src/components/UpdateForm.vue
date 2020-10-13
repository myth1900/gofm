<template>
  <v-form>
    <v-container>
      <v-row>
        <v-col
        cols="12"
        md="5"
        >
          <v-text-field
              v-model="room_id"
              label="房间ID"
              required
          >
          </v-text-field>
        </v-col>

        <v-col
            cols="12"
            md="5"
        >
          <v-text-field
              v-model="nums"
              label="人数"
              required
          >
          </v-text-field>
        </v-col>


        <v-col
            cols="12"
            md="2"
        >
          <v-btn class="mx-2" fab dark color="indigo"
                 v-on:click="submitUpdate"
          >
            <v-icon dark>mdi-check-bold</v-icon>
          </v-btn>
        </v-col>


      </v-row>
    </v-container>
  </v-form>
</template>

<script>
export default {
name: "UpdateForm",
  data: function (){
  return {
    room_id: '',
    nums: '',
  }
  },
  methods: {
  submitUpdate() {
    let url = "/api" + "/room/" + this.room_id + "/audience"
    this.$axios({
      url: url,
      method: "put",
      params: {
        nums: this.nums,
      }
    }).then(res => {
      console.log(res)
    })
  }
  }
}
</script>

<style scoped>

</style>