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
                    <span class="caption-subject font-green sbold">Actions  </span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                 <!--<div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Operation: </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>-->
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{range $index,$act := .Actions}}
                                  <tr><td>ActionID:</td><td>{{$act.ActID}}</td></tr>
                                  <tr><td>TxHash:</td><td>{{$act.TxHash}}</td></tr>
                                  <tr><td>Src:</td><td>{{$act.Src}}</td></tr>
                                  <tr><td>Dst:</td><td>{{$act.Dst}}</td></tr>
                                  <tr><td>Amount:</td><td>{{$act.Amount}}</td></tr>
                                  <tr><td>Body:</td><td><textarea style="width:80%;height:50px;" readonly>{{$act.Body}}</textarea></td></tr>
                                  <tr><td>Memo:</td><td>{{$act.Memo}}</td></tr>
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