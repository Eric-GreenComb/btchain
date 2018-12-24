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
                    <span class="caption-subject font-green sbold">Transactions  </span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                 <!--<div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Transactions </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>-->
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{range $index,$op := .TransDetail}}
								<tr><td><span class="caption-subject font-green sbold">Action Number:</span></td><td><span class="caption-subject font-green sbold">{{$index}}</span></td></tr>
								<tr><td>TxHash:</td><td><a href=/view/trans/detail/{{$op.TxHash}}>{{$op.TxHash}}</a></td></tr>
								<tr><td>Height:</td><td>{{$op.BlockHeight}}</td></tr>
								<tr><td>BlockHash:</td><td><a href=/view/blocks/hash/{{$op.BlockHash}}>{{$op.BlockHash}}</a></td></tr>
								<tr><td>Actions:</td><td>{{$op.ActionCount}}</td></tr>
								<tr><td>ActionID:</td><td>{{$op.ActionID}}</td></tr>
								<tr><td>Src:</td><td>{{$op.Src}} <a href=/view/accounts/{{$op.Src}}/payout>payout</a>:<a href=/view/accounts/{{$op.Src}}/income>income</a></td></tr>
								<tr><td>Nonce:</td><td>{{$op.Nonce}}</td></tr>
								<tr><td>Dst:</td><td>{{$op.Dst}} <a href=/view/accounts/{{$op.Dst}}/payout>payout</a>:<a href=/view/accounts/{{$op.Dst}}/income>income</a></td></tr>
								<tr><td>Amount:</td><td>{{$op.Amount}}</td></tr>
								<tr><td>Data:</td><td><textarea style="width:100%;height:40px;" readonly>{{$op.JData}}</textarea></td></tr>
								<tr><td>Memo:</td><td>{{$op.Memo}}</td></tr>
								<tr><td>Time:</td><td>{{$op.CreateAt}}</td></tr>
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