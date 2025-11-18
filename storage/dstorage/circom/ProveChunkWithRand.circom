pragma circom 2.2.2;
include "poseidon.circom";
include "comparators.circom";

template ProveChunkWithRand {
    signal input elems[133];
    signal input rand[2];
    // 公共输入
    signal input expected_chunk_hash;
    signal input expected_hash;    // Poseidon 输出 field

    signal output is_valid;

    // 前面只计算文件分片的poseidon hash, 最后二合一的时候，合入随机数，这样计算总的poseidon hash
    // 计算第一层： 133个节点 67 = (133 + 1) / 2
    component h0[67];
    signal layer0[67];
    for (var i = 0; i < 67; i++) {
        h0[i] = Poseidon(3);
        h0[i].inputs[0] <== elems[2*i];
        // 如果是奇数，补 0
        h0[i].inputs[1] <== (2*i+1 < 133 ? elems[2*i+1] : 0);
        h0[i].inputs[2] <== 0;
        layer0[i] <== h0[i].out;
    }

    // 计算第二层： 67个节点
    component h1[34];
    signal layer1[34];
    for (var i = 0; i < 34; i++) {
        h1[i] = Poseidon(3);
        h1[i].inputs[0] <== layer0[2*i];
        // 如果是奇数，补 0
        h1[i].inputs[1] <== (2*i+1 < 67 ? layer0[2*i+1] : 0);
        h1[i].inputs[2] <== 0;
        layer1[i] <== h1[i].out;
    }

    // 计算第三层： 34个节点
    component h2[17];
    signal layer2[17];
    for (var i = 0; i < 17; i++) {
        h2[i] = Poseidon(3);
        h2[i].inputs[0] <== layer1[2*i];
        h2[i].inputs[1] <== layer1[2*i+1];
        h2[i].inputs[2] <== 0;
        layer2[i] <== h2[i].out;
    }

    // 计算第四层： 17个节点
    component h3[9];
    signal layer3[9];
    for (var i = 0; i < 9; i++) {
        h3[i] = Poseidon(3);
        h3[i].inputs[0] <== layer2[2*i];
        h3[i].inputs[1] <== (2*i+1 < 17 ? layer2[2*i+1] : 0);
        h3[i].inputs[2] <== 0;
        layer3[i] <== h3[i].out;
    }

    // 计算第五层： 9个节点
    component h4[5];
    signal layer4[5];
    for (var i = 0; i < 5; i++) {
        h4[i] = Poseidon(3);
        h4[i].inputs[0] <== layer3[2*i];
        // 奇数，补 0
        h4[i].inputs[1] <== (2*i+1 < 9 ? layer3[2*i+1] : 0);
        h4[i].inputs[2] <== 0;
        layer4[i] <== h4[i].out;
    }

    // 计算第六层： 5个节点
    component h5[3];
    signal layer5[3];
    for (var i = 0; i < 3; i++) {
        h5[i] = Poseidon(3);
        h5[i].inputs[0] <== layer4[2*i];
        // 如果是奇数，补 0
        h5[i].inputs[1] <== (2*i+1 < 5 ? layer4[2*i+1] : 0);
        h5[i].inputs[2] <== 0;
        layer5[i] <== h5[i].out;
    }

    // 计算第七层： 2个节点
    component h6[2];
    signal layer6[2];
    for (var i = 0; i < 2; i++) {
        h6[i] = Poseidon(3);
        h6[i].inputs[0] <== layer5[2*i];
        h6[i].inputs[1] <== (2*i+1 < 3 ? layer5[2*i+1] : 0);
        h6[i].inputs[2] <== 0;
        layer6[i] <== h6[i].out;
    }

    // 计算第八层： 1个节点
    component h7[1];
    signal layer7[1];
    h7[0] = Poseidon(3);
    h7[0].inputs[0] <== layer6[0];
    h7[0].inputs[1] <== layer6[1];
    h7[0].inputs[2] <== 0;
    layer7[0] <== h7[0].out;

    component chunk_eq = IsEqual();
    chunk_eq.in[0] <== layer7[0];
    chunk_eq.in[1] <== expected_chunk_hash;

    component h_rand[1];
    signal layer_rand[1];
    h_rand[0] = Poseidon(3);
    h_rand[0].inputs[0] <== rand[0];
    h_rand[0].inputs[1] <== rand[1];
    h_rand[0].inputs[2] <== 0;
    layer_rand[0] <== h_rand[0].out;

    component h8[1];
    signal layer8[1];
    h8[0] = Poseidon(3);
    h8[0].inputs[0] <== layer7[0];
    h8[0].inputs[1] <== layer_rand[0];
    h8[0].inputs[2] <== 0;
    layer8[0] <== h8[0].out;

    component eq = IsEqual();
    eq.in[0] <== layer8[0];
    eq.in[1] <== expected_hash;

    is_valid <== eq.out * chunk_eq.out;
}

component main {public [expected_chunk_hash, expected_hash]} = ProveChunkWithRand();