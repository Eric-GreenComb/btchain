<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
    }
</style>

<div class="row">
    <div class="col-md-12">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-title">
                <div class="caption">
                    <i class="fa fa-book font-green"></i>
                    <span class="caption-subject font-green sbold">Hash : {{.Block.Hash}}</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Chain ID:</td><td>{{.Block.ChainID}}</td></tr>
                          <tr><td>Height:</td><td>{{.Block.Height}}</td></tr>
                          <tr><td>NumTxs:</td><td>{{.Block.NumTx}}</td></tr>
                          <tr><td>ValidatorsHash:</td><td>{{.Block.ValidHash}}</td></tr>
                          <tr><td>AppHash:</td><td>{{.Block.AppHash}}</td></tr>
                          <tr><td>Time:</td><td>{{.Block.Time}}</td></tr>
                        </tbody>
                    </table>
                </div>

                 <div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Transactions : </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{range $index,$tx := .Transactions}}
                                  <tr><td>Hash:</td><td>{{$tx.Hash}}</td></tr>
                                  <tr><td>BaseFee:</td><td>{{$tx.BaseFee}}</td></tr>
                                  <tr><td>TxType:</td><td>{{$tx.TxType}}</td></tr>
								  <tr><td>Actions:</td><td>{{$tx.Actions}}</td></tr>
								  <tr><td title={{$tx.TxID}} >Actions:</td><td><a href=/view/trans/txid/{{$tx.TxID}}> View >></a></td></tr>
                             {{end}}
                            </tbody>
                        </table>
                    </div>
                </div>

            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>

<script>
	//var hrefOnClick = function(){
	//	window.location.href='/view/trans/txid/'+;
	//}
</script>